from fastapi import APIRouter
from app.core.config import INSTANCE_NAME
from app.core.logger import logger
from app.models.schemas import PurchaseEvent, OrderResponseDTO
from app.mappings.purchase_mapper import map_purchase_event_to_order_payload
from app.services.orders_service import dispatch_order_to_service, get_order_from_service
from app.services.items_service import get_item_from_service

router = APIRouter()

@router.get("")
def gateway_root():
    return {
        "status": "success",
        "message": "Response from the Gateway!",
        "instance": INSTANCE_NAME,
        "path": "/gateway"
    }

@router.post("/orders", response_model=OrderResponseDTO)
def receive_purchase_event(event: PurchaseEvent):
    logger.info(f"[{INSTANCE_NAME}] === NEW PURCHASE EVENT RECEIVED ===")
    logger.info(f"[{INSTANCE_NAME}] User: {event.user} | Product: {event.item_id} | Quantity: {event.quantity} | Total: ${event.amount:,.2f} ARS")
    
    # 1. Map incoming event to the payload required by orders-service
    order_payload = map_purchase_event_to_order_payload(event)
    
    # 2. Dispatch mapped payload using the service layer
    result = dispatch_order_to_service(order_payload)
    
    return result

@router.get("/orders/{order_id}")
def get_order(order_id: str):
    return get_order_from_service(order_id)

@router.get("/items/{item_id}")
def get_item(item_id: str):
    return get_item_from_service(item_id)

@router.get("/info")
def gateway_info():
    return {
        "instance": INSTANCE_NAME,
        "details": "This is the specific path routed by Nginx using Round Robin."
    }
