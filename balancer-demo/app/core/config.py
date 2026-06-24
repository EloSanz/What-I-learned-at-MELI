import os

INSTANCE_NAME = os.getenv("INSTANCE_NAME", "Unknown")
ORDERS_SERVICE_URL = os.getenv("ORDERS_SERVICE_URL", "http://orders-service:8082")
ITEMS_SERVICE_URL = os.getenv("ITEMS_SERVICE_URL", "http://items-service:8081")