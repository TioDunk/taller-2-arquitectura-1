from fastapi import FastAPI as FastAPIB
from pydantic import BaseModel as BaseModelB, Field as FieldB
from enum import Enum as EnumB
from typing import Dict as DictB
import os as osb
import requests as rq
from datetime import datetime as dt
import logging as logging_b
try:
    from colorama import Fore as ForeB, Style as StyleB
except Exception:
    class _FB: MAGENTA=""; RED=""
    class _SB: RESET_ALL=""
    ForeB = _FB(); StyleB = _SB()

# ---- Logging (Billing) ----
BLOG = logging_b.getLogger("BILLING")
BLOG.setLevel(logging_b.INFO)
if not logging_b.getLogger().handlers:
    logging_b.basicConfig(level=logging_b.INFO, format="%(asctime)s | %(levelname)s | %(name)s | %(message)s")

def blog_info(msg: str):
    BLOG.info(ForeB.MAGENTA + "[BILLING] " + StyleB.RESET_ALL + msg)

# --------- Dominio (Billing) ---------
class PaymentStatus(str, EnumB):
    APPROVED = "APPROVED"
    REJECTED = "REJECTED"

class MoneyB(BaseModelB):
    amount_cents: int = FieldB(ge=0)
    currency: str = "COP"

class Payment(BaseModelB):
    id: str
    order_id: str
    amount: MoneyB
    status: PaymentStatus
    processed_at: dt = FieldB(default_factory=dt.utcnow)

# --------- Repositorio (in-memory) ---------
class PaymentRepository:
    def __init__(self) -> None:
        self._db: DictB[str, Payment] = {}

    def save(self, p: Payment) -> None:
        self._db[p.id] = p

    def get(self, pid: str) -> Payment:
        return self._db[pid]

# --------- Servicio de aplicación (Billing) ---------
class BillingService:
    def __init__(self, repo: PaymentRepository) -> None:
        self.repo = repo

    def register_payment(self, order_id: str, amount: MoneyB) -> Payment:
        # PoC: aprobamos siempre. Aquí podrías validar con pasarela/fraude/etc.
        pid = f"PAY-{order_id}"
        payment = Payment(id=pid, order_id=order_id, amount=amount, status=PaymentStatus.APPROVED)
        self.repo.save(payment)
        return payment

# --------- FastAPI + publicación de evento a Sales ---------
app = FastAPIB(title="Billing Service (DDD PoC)")
repo_p = PaymentRepository()
billing = BillingService(repo_p)
SALES_URL = osb.getenv("SALES_URL", "http://localhost:8000")

class CreatePaymentDTO(BaseModelB):
    order_id: str
    amount_cents: int = FieldB(ge=0)
    currency: str = "COP"

@app.get("/health")
def health_b():
    return {"service": "billing", "status": "ok"}

@app.post("/payments")
def create_payment(dto: CreatePaymentDTO):
    blog_info(f"⇢ POST /payments recibido para orden {dto.order_id} monto={dto.amount_cents} {dto.currency}")
    p = billing.register_payment(order_id=dto.order_id, amount=MoneyB(amount_cents=dto.amount_cents, currency=dto.currency))
    blog_info(f"Pago {p.id} APROBADO. Notificando a Sales {SALES_URL}/integration/payment-approved")
    evt = {"order_id": p.order_id, "amount_cents": p.amount.amount_cents, "currency": p.amount.currency}
    try:
        resp = rq.post(f"{SALES_URL}/integration/payment-approved", json=evt, timeout=5)
        blog_info(f"Respuesta de Sales {resp.status_code}: {resp.text[:120]}")
    except Exception as e:
        BLOG.error(ForeB.RED + f"Error notificando a Sales: {e}")
    return {"payment_id": p.id, "status": p.status}

@app.get("/payments/{payment_id}")
def get_payment(payment_id: str):
    p = repo_p.get(payment_id)
    return p

