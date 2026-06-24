from app.models.schemas import PurchaseEvent, OrderPayloadDTO

def map_purchase_event_to_order_payload(event: PurchaseEvent) -> OrderPayloadDTO:
    """
    Transforms the incoming PurchaseEvent into the DTO expected
    by the orders-service.
    """
    return OrderPayloadDTO(
        user_id=event.user,
        item_id=event.item_id,
        quantity=event.quantity,
        amount=event.amount,
        address=event.address.get_secret_value()
    )
