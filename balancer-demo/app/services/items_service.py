import json
from urllib.request import Request, urlopen
from app.core.config import ITEMS_SERVICE_URL

def get_item_from_service(item_id: str) -> dict:
    """
    Communicates with the internal items-service to fetch item details.
    """
    url = f"{ITEMS_SERVICE_URL}/api/items/{item_id}"
    try:
        req = Request(url, method="GET")
        with urlopen(req, timeout=5) as response:
            resp_body = response.read().decode("utf-8")
            resp_json = json.loads(resp_body)
            return resp_json
    except Exception as e:
        return {
            "error": "Items service offline",
            "message": f"Could not contact items service: {str(e)}"
        }
