from fastapi import Request
from app.core.logger import logger
from app.core.config import INSTANCE_NAME

# Middleware to intercept and log incoming connections
async def log_requests(request: Request, call_next):
    client = f"{request.client.host}:{request.client.port}" if request.client else "unknown"
    logger.info(f"[{INSTANCE_NAME}] New connection accepted from {client} -> {request.method} {request.url.path}")
    response = await call_next(request)
    return response
