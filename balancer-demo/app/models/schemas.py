from pydantic import BaseModel, Field, SecretStr

class PurchaseEvent(BaseModel):
    """
    PurchaseEvent represents the payload received from the Next.js frontend.
    
    Flow:
    1. Originates from the client browser (Next.js checkout UI).
    2. Enters through the Nginx Load Balancer (port 8080).
    3. Received by one of the FastAPI gateway instances (port 8000).
    4. Data is parsed, validated, and forwarded to the internal Orders Service (port 8082).
    """
    event: str = Field(..., description="The type of event, e.g., 'buy_now'")
    user: str = Field(..., description="ID of the user making the purchase")
    item: str = Field(..., description="ID of the product being purchased")
    quantity: int = Field(..., description="Number of units bought", gt=0)
    address: SecretStr = Field(..., description="User's shipping address. Obfuscated in logs using SecretStr for privacy.")
    amount: float = Field(..., description="Total price of the purchase in local currency")
    timestamp: str = Field(..., description="ISO 8601 timestamp of the transaction")
