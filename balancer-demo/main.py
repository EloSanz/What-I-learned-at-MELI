import os
import json
import urllib.request
from fastapi import FastAPI, Request
from pydantic import BaseModel

app = FastAPI()

import logging

# Configuración del logger
logging.basicConfig(level=logging.INFO, format="%(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

INSTANCE_NAME = os.getenv("INSTANCE_NAME", "Desconocida")

# Modelo de datos para recibir el evento de compra
class PurchaseEvent(BaseModel):
    event: str
    user: str
    item: str
    quantity: int
    address: str
    amount: float
    timestamp: str

# Middleware para interceptar e imprimir las conexiones entrantes por consola
@app.middleware("http")
async def log_requests(request: Request, call_next):
    client = f"{request.client.host}:{request.client.port}" if request.client else "desconocido"
    logger.info(f"[{INSTANCE_NAME}] Nueva conexión aceptada desde {client} -> {request.method} {request.url.path}")
    response = await call_next(request)
    return response

@app.get("/gateway")
def gateway_root():
    return {
        "status": "success",
        "message": "¡Respuesta desde el Gateway!",
        "instancia": INSTANCE_NAME,
        "path": "/gateway"
    }

@app.post("/gateway")
def receive_purchase_event(event: PurchaseEvent):
    logger.info(f"[{INSTANCE_NAME}] === NUEVO EVENTO DE COMPRA RECIBIDO ===")
    logger.info(f"[{INSTANCE_NAME}] Usuario: {event.user} | Producto: {event.item} | Cantidad: {event.quantity} | Total: ${event.amount:,.2f} ARS")
    
    # 1. Preparar datos para enviar a orders-service
    orders_service_url = os.getenv("ORDERS_SERVICE_URL", "http://orders-service:8082")
    url = f"{orders_service_url}/api/orders"
    
    order_payload = {
        "user_id": event.user,
        "item_id": "MLA43960787",  # ID del monitor gamer Xiaomi sembrado en la base de datos
        "quantity": event.quantity,
        "amount": event.amount,
        "address": event.address
    }
    
    headers = {"Content-Type": "application/json"}
    
    try:
        req = urllib.request.Request(
            url,
            data=json.dumps(order_payload).encode("utf-8"),
            headers=headers,
            method="POST"
        )
        # Timeout de 5 segundos
        with urllib.request.urlopen(req, timeout=5) as response:
            resp_body = response.read().decode("utf-8")
            resp_json = json.loads(resp_body)
            
            if resp_json.get("status") == "success":
                order_data = resp_json.get("data")
                logger.info(f"[{INSTANCE_NAME}] Orden creada correctamente en orders-service (ID: {order_data.get('id')})")
                return {
                    "status": "event_received",
                    "processed_by": INSTANCE_NAME,
                    "event_id": order_data.get("id"),
                    "message": f"Orden creada en orders-service con estado: {order_data.get('status')}"
                }
            else:
                return {
                    "error": "Orders service validation failed",
                    "message": resp_json.get("message", "Error desconocido")
                }
    except Exception as e:
        logger.error(f"[{INSTANCE_NAME}] Error al despachar orden a orders-service: {e}")
        return {
            "error": "Orders service offline",
            "message": f"No se pudo contactar al servicio de órdenes: {str(e)}"
        }

@app.get("/gateway/info")
def gateway_info():
    return {
        "instancia": INSTANCE_NAME,
        "detalles": "Este es el path específico redirigido por Nginx usando Round Robin."
    }


