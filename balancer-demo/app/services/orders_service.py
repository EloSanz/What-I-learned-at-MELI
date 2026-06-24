import json
from urllib.request import Request, urlopen
from urllib.error import HTTPError
from app.core.config import ORDERS_SERVICE_URL, INSTANCE_NAME
from app.core.logger import logger
from app.models.schemas import OrderPayloadDTO, OrderResponseDTO

def dispatch_order_to_service(order_payload: OrderPayloadDTO) -> OrderResponseDTO:
    """
    Sends the mapped order payload to the internal orders-service via HTTP POST.
    """
    url = f"{ORDERS_SERVICE_URL}/api/orders"
    headers = {"Content-Type": "application/json"}
    
    try:
        req = Request(
            url,
            data=order_payload.model_dump_json().encode("utf-8"),
            headers=headers,
            method="POST"
        )
        # 5-second timeout
        with urlopen(req, timeout=5) as response:
            resp_body = response.read().decode("utf-8")
            resp_json = json.loads(resp_body)
            
            if resp_json.get("status") == "success":
                order_data = resp_json.get("data")
                logger.info(f"[{INSTANCE_NAME}] Order created successfully in orders-service (ID: {order_data.get('id')})")
                return OrderResponseDTO(
                    status="event_received",
                    processed_by=INSTANCE_NAME,
                    event_id=order_data.get("id"),
                    message=f"Order created in orders-service with status: {order_data.get('status')}"
                )
            else:
                return OrderResponseDTO(
                    error="Orders service validation failed",
                    message=resp_json.get("message", "Unknown error")
                )
    except Exception as e:
        logger.error(f"[{INSTANCE_NAME}] Error dispatching order to orders-service: {e}")
        return OrderResponseDTO(
            error="Orders service offline",
            message=f"Could not contact orders service: {str(e)}"
        )

def get_order_from_service(order_id: str) -> dict:
    url = f"{ORDERS_SERVICE_URL}/api/orders/{order_id}"
    
    try:
        req = Request(url, method="GET")
        with urlopen(req, timeout=5) as response:
            resp_body = response.read().decode("utf-8")
            return json.loads(resp_body)
    except HTTPError as e:
        if e.code == 404:
            return {"error": "Order not found"}
        return {"error": f"Orders service error: {e.reason}"}
    except Exception as e:
        logger.error(f"[{INSTANCE_NAME}] Error fetching order from orders-service: {e}")
        return {"error": "Could not contact orders service"}
