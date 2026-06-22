from fastapi import FastAPI
from starlette.middleware.base import BaseHTTPMiddleware

from app.api.routes import gateway
from app.middlewares.logging import log_requests

app = FastAPI(title="Meli Gateway Demo")

# Add Middlewares
app.add_middleware(BaseHTTPMiddleware, dispatch=log_requests)

# Include Routers
app.include_router(gateway.router, prefix="/gateway", tags=["gateway"])
