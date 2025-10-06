from __future__ import annotations
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from enum import Enum
from typing import List, Dict, Optional
from datetime import datetime
import os
import requests
import logging

try:
    from colorama import init as _color_init, Fore, Style
    _color_init(autoreset=True)
except Exception:
    class _F: CYAN = RED = ""
    class _S: RESET_ALL = ""
    Fore = _F(); Style = _S()

# ---- Logging (Sales) ----
SLOG = logging.getLogger("SALES")
SLOG.setLevel(logging.INFO)
if not logging.getLogger().handlers:
    logging.basicConfig(level=logging.INFO, format="%(asctime)s | %(levelname)s | %(name)s | %(message)s")

def slog_info(msg: str):
    SLOG.info(Fore.CYAN + "[SALES] " + Style.RESET_ALL + msg)

# --------- Dominio (Sales) ---------
class DomainError(Exception):
    pass

class OrderStatus(str, Enum):
    PENDING = "PENDING"
    PAID = "PAID"

class Money(BaseModel):
    amount_cents: int = Field(ge=0)
    currency: str = "COP"

    def __add__(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise DomainError("Monedas distintas")
        return Money(amount_cents=self.amount_cents + other.amount_cents, currency=self.currency)

class OrderItem(BaseModel):
    sku: str
    qty: int = Field(ge=1)
    unit_price: Money

    @property
    def line_total(self) -> Money:
        return Money(amount_cents=self.unit_price.amount_cents * self.qty, currency=self.unit_price.currency)

class Order(BaseModel):
    id: str
    customer_id: str
    items: List[OrderItem]
    status: OrderStatus = OrderStatus.PENDING
    created_at: datetime = Field(default_factory=datetime.utcnow)

    def total(self) -> Money:
        total = Money(amount_cents=0, currency=self.items[0].unit_price.currency if self.items else "COP")
        for it in self.items:
            total = total + it.line_total
        return total

    def mark_paid(self) -> None:
        if self.status == OrderStatus.PAID:
            return
        if self.status != OrderStatus.PENDING:
            raise DomainError("Sólo se pueden pagar pedidos PENDING")
        self.status = OrderStatus.PAID

# --------- Repositorio (in-memory) ---------
class OrderRepository:
    def __init__(self) -> None:
        self._db: Dict[str, Order] = {}

    def save(self, order: Order) -> None:
        self._db[order.id] = order

    def get(self, order_id: str) -> Order:
        if order_id not in self._db:
            raise DomainError("Pedido no encontrado")
        return self._db[order_id]

# --------- Servicio de aplicación (Sales) ---------
class SalesService:
    def __init__(self, repo: OrderRepository) -> None:
        self.repo = repo

    def place_order(self, order_id: str, customer_id: str, items: List[OrderItem]) -> Order:
        if not items:
            raise DomainError("Pedido requiere al menos un ítem")
        order = Order(id=order_id, customer_id=customer_id, items=items)
        self.repo.save(order)
        return order

    def mark_order_as_paid(self, order_id: str) -> Order:
        order = self.repo.get(order_id)
        order.mark_paid()
        self.repo.save(order)
        return order

# --------- FastAPI + publicación de evento a Billing ---------
app = FastAPI(title="Sales Service (DDD PoC)")
repo = OrderRepository()
sales = SalesService(repo)
BILLING_URL = os.getenv("BILLING_URL", "http://localhost:8001")

class CreateOrderItemDTO(BaseModel):
    sku: str
    qty: int = Field(ge=1)
    unit_price_cents: int = Field(ge=0)
    currency: str = "COP"

class CreateOrderDTO(BaseModel):
    order_id: str
    customer_id: str
    items: List[CreateOrderItemDTO]

class OrderVM(BaseModel):
    id: str
    status: OrderStatus
    total_amount_cents: int
    currency: str

# Healthcheck (observabilidad / readiness)
@app.get("/health")
def health():
    return {"service": "sales", "status": "ok"}

# Caso de uso: crear pedido + publicar OrderPlaced hacia Billing (vía HTTP)
@app.post("/orders", response_model=OrderVM)
def create_order(dto: CreateOrderDTO):
    slog_info(f"⇢ POST /orders recibido para {dto.order_id}")
    items = [OrderItem(sku=i.sku, qty=i.qty, unit_price=Money(amount_cents=i.unit_price_cents, currency=i.currency)) for i in dto.items]
    order = sales.place_order(dto.order_id, dto.customer_id, items)
    total = order.total()
    slog_info(f"Pedido {order.id} total {total.amount_cents} {total.currency}")
    payload = {"order_id": order.id, "amount_cents": total.amount_cents, "currency": total.currency}
    try:
        resp = requests.post(f"{BILLING_URL}/payments", json=payload, timeout=5)
        slog_info(f"→ Enviado a Billing ({resp.status_code})")
    except Exception as e:
        print("Error al contactar Billing:", e)
    return OrderVM(id=order.id, status=order.status, total_amount_cents=total.amount_cents, currency=total.currency)

# Consulta: obtener estado del pedido
@app.get("/orders/{order_id}", response_model=OrderVM)
def get_order(order_id: str):
    try:
        order = repo.get(order_id)
        total = order.total()
        return OrderVM(id=order.id, status=order.status, total_amount_cents=total.amount_cents, currency=total.currency)
    except DomainError as e:
        # 404 cuando el pedido no existe
        raise HTTPException(status_code=404, detail=str(e))

class PaymentApprovedDTO(BaseModel):
    order_id: str
    amount_cents: int
    currency: str

# Webhook: evento de Billing → Sales (PaymentApproved)
@app.post("/integration/payment-approved")
def on_payment_approved(evt: PaymentApprovedDTO):
    slog_info(f"↩ Evento PaymentApproved para orden {evt.order_id}")
    order = repo.get(evt.order_id)
    if evt.amount_cents >= order.total().amount_cents:
        sales.mark_order_as_paid(order.id)
        return {"ok": True, "status": order.status}
    else:
        return {"ok": False, "reason": "Pago insuficiente"}
