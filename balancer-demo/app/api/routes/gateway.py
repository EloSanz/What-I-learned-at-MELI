import json
import urllib.request
from fastapi import APIRouter
from app.core.config import INSTANCE_NAME, ORDERS_SERVICE_URL
from app.core.logger import logger
from app.models.schemas import PurchaseEvent

router = APIRouter()

@router.get("")
def gateway_root():
    return {
        "status": "success",
        "message": "Response from the Gateway!",
        "instance": INSTANCE_NAME,
        "path": "/gateway"
    }

@router.post("")
def receive_purchase_event(event: PurchaseEvent):
    logger.info(f"[{INSTANCE_NAME}] === NEW PURCHASE EVENT RECEIVED ===")
    logger.info(f"[{INSTANCE_NAME}] User: {event.user} | Product: {event.item} | Quantity: {event.quantity} | Total: ${event.amount:,.2f} ARS")
    
    # 1. Prepare data for orders-service
    url = f"{ORDERS_SERVICE_URL}/api/orders"
    
    order_payload = {
        "user_id": event.user,
        "item_id": "MLA43960787",  # Seeded Xiaomi gamer monitor ID
        "quantity": event.quantity,
        "amount": event.amount,
        "address": event.address.get_secret_value()
    }
    
    headers = {"Content-Type": "application/json"}
    
    try:
        req = urllib.request.Request(
            url,
            data=json.dumps(order_payload).encode("utf-8"),
            headers=headers,
            method="POST"
        )
        # 5-second timeout
        with urllib.request.urlopen(req, timeout=5) as response:
            resp_body = response.read().decode("utf-8")
            resp_json = json.loads(resp_body)
            
            if resp_json.get("status") == "success":
                order_data = resp_json.get("data")
                logger.info(f"[{INSTANCE_NAME}] Order created successfully in orders-service (ID: {order_data.get('id')})")
                return {
                    "status": "event_received",
                    "processed_by": INSTANCE_NAME,
                    "event_id": order_data.get("id"),
                    "message": f"Order created in orders-service with status: {order_data.get('status')}"
                }
            else:
                return {
                    "error": "Orders service validation failed",
                    "message": resp_json.get("message", "Unknown error")
                }
    except Exception as e:
        logger.error(f"[{INSTANCE_NAME}] Error dispatching order to orders-service: {e}")
        return {
            "error": "Orders service offline",
            "message": f"Could not contact orders service: {str(e)}"
        }

@router.get("/info")
def gateway_info():
    return {
        "instance": INSTANCE_NAME,
        "details": "This is the specific path routed by Nginx using Round Robin."
    }
