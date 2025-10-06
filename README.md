# 🧩 Prueba de Concepto (PoC) — Domain-Driven Design con Microservicios

# 🎦 Video demostrativo:
https://www.loom.com/share/2747252d333749b6a070251745fa8dd2

Este proyecto demuestra cómo aplicar **Domain-Driven Design (DDD)** en una arquitectura basada en **microservicios** utilizando **FastAPI** (Python).
La PoC está compuesta por **dos microservicios** independientes que simulan un flujo de negocio típico: **pedidos (Sales)** y **pagos (Billing)**.

---

## 🧠 Objetivo

Mostrar cómo el diseño centrado en el dominio permite:
- Separar responsabilidades en microservicios (**granularidad adecuada**).
- Definir **contextos delimitados (bounded contexts)**.
- Comunicar servicios mediante **eventos de dominio**.
- Reflejar las reglas de negocio de forma clara y desacoplada.

---

## ⚙️ Arquitectura general

```
┌────────────────────────┐          ┌───────────────────────────┐
│      Sales Service      │          │       Billing Service      │
│────────────────────────│          │────────────────────────────│
│ Crea pedidos (Order)    │ ───────▶ │ Registra pagos (Payment)   │
│ Publica “OrderPlaced”   │          │ Publica “PaymentApproved”  │
│ Actualiza estado (PAID) │ ◀─────── │ Notifica a Sales           │
└────────────────────────┘          └───────────────────────────┘
```

---

## 🧱 Servicios implementados

### **1️⃣ Sales Service (`sales_service.py`)**
Encargado del ciclo de vida de los **pedidos**.

**Endpoints:**
| Método | Ruta | Descripción |
|---------|------|-------------|
| `GET` | `/health` | Verifica que el servicio esté corriendo |
| `POST` | `/orders` | Crea un pedido y envía evento a Billing |
| `GET` | `/orders/{order_id}` | Consulta el estado del pedido |
| `POST` | `/integration/payment-approved` | Recibe evento `PaymentApproved` desde Billing |

---

### **2️⃣ Billing Service (`billing_service.py`)**
Encargado de los **pagos** asociados a los pedidos.

**Endpoints:**
| Método | Ruta | Descripción |
|---------|------|-------------|
| `GET` | `/health` | Verifica que el servicio esté corriendo |
| `POST` | `/payments` | Registra un pago y notifica a Sales |
| `GET` | `/payments/{payment_id}` | Consulta un pago registrado |

---

## 🚀 Ejecución local

### 🔧 1. Clonar e instalar dependencias
```bash
git clone https://github.com/tu-usuario/taller-2-arquitectura-1.git
cd taller-2-arquitectura-1
pip install fastapi uvicorn requests colorama
```

### 🖥️ 2. Ejecutar los microservicios

#### Terminal 1 — **Sales Service**
```bash
uvicorn sales_service:app --reload --port 8000
```

#### Terminal 2 — **Billing Service**
```bash
uvicorn billing_service:app --reload --port 8001
```

---

## 🧪 3. Prueba del flujo completo

### ➤ Crear pedido
```bash
curl -X POST http://localhost:8000/orders \
-H "Content-Type: application/json" \
-d '{
  "order_id": "ORD-1001",
  "customer_id": "CUST-77",
  "items": [
    {"sku": "SKU-1", "qty": 2, "unit_price_cents": 115000}
  ]
}'
```

📜 Resultado esperado:
- **Sales** registra el pedido.
- **Billing** recibe el evento y aprueba el pago.
- **Sales** marca el pedido como pagado (`status: PAID`).

---

### ➤ Consultar pedido
```bash
curl http://localhost:8000/orders/ORD-1001
```
📦 Resultado:
```json
{
  "id": "ORD-1001",
  "status": "PAID",
  "total_amount_cents": 230000,
  "currency": "COP"
}
```

---

## 🧩 Conceptos de DDD aplicados

| Concepto | Aplicación en la PoC |
|-----------|----------------------|
| **Granularidad adecuada** | Cada microservicio tiene una sola responsabilidad. |
| **Bounded Contexts** | `Sales` y `Billing` definen dominios independientes. |
| **Lenguaje ubicuo** | Los nombres reflejan conceptos del negocio (pedido, pago, total). |
| **Eventos de dominio** | `OrderPlaced` y `PaymentApproved` comunican los contextos. |
| **Autonomía y bajo acoplamiento** | Comunicación vía HTTP, sin compartir base de datos. |
| **Event Storming (implícito)** | El flujo de eventos guía la definición de los contextos. |

---

## 🧾 Logs de ejemplo

```
[SALES] ⇢ POST /orders recibido para ORD-1001
[SALES] Pedido ORD-1001 total 230000 COP
[BILLING] ⇢ POST /payments recibido para orden ORD-1001 monto=230000 COP
[BILLING] Pago PAY-ORD-1001 APROBADO. Notificando a Sales
[SALES] ↩ Evento PaymentApproved para orden ORD-1001
```

---

## 👥 Autores

- **Grupo 13**
- **Grupo G1 — Diplomado en Arquitectura de Software**
- Universidad de La Sabana – 2025

---

📚 **Tecnologías:** Python · FastAPI · Uvicorn
🧩 **Arquitectura:** Microservicios · DDD · Event Storming (teórico)
