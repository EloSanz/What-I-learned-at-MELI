import os
import logging
import logging_loki
from fastmcp import FastMCP
from loki_client import query_loki

LOKI_URL = os.getenv("LOKI_URL", "http://localhost:3100")
loki_handler = logging_loki.LokiHandler(
    url=f"{LOKI_URL}/loki/api/v1/push",
    tags={"service": "meli-logs-server"},
    version="1",
)
logging.basicConfig(level=logging.INFO, handlers=[logging.StreamHandler(), loki_handler])
logger = logging.getLogger("meli-logs-server")

mcp = FastMCP("Meli Logs Server 📝")

@mcp.prompt()
def available_services() -> str:
    """Get the list of available services to query logs for."""
    services = ["gateway-a", "orders-service", "items-service", "meli-logs-server"]
    return f"The following services are available for log analysis: {', '.join(services)}. Please use the get_logs or search_logs tools to fetch their logs."

@mcp.tool
def get_logs(service_name: str, start_time: str, end_time: str, limit: int = 100) -> str:
    """
    Fetch logs for a specific service within a time range.
    
    Args:
        service_name: The name of the docker-compose service (e.g., 'gateway-a', 'orders-service', 'items-service').
        start_time: ISO 8601 start time (e.g., '2026-06-24T00:00:00Z').
        end_time: ISO 8601 end time (e.g., '2026-06-25T00:00:00Z').
        limit: Maximum number of log lines to return (default 100).
    """
    logger.info(f"Executing get_logs for service={service_name}")
    query = f'{{service="{service_name}"}}'
    return query_loki(query, start_time, end_time, limit)

@mcp.tool
def search_logs(service_name: str, search_query: str, start_time: str, end_time: str, limit: int = 100) -> str:
    """
    Search for a specific text string within the logs of a service.
    
    Args:
        service_name: The name of the docker-compose service.
        search_query: The text string to search for (e.g., 'ERROR', 'UUID').
        start_time: ISO 8601 start time.
        end_time: ISO 8601 end time.
        limit: Maximum number of log lines to return.
    """
    logger.info(f"Executing search_logs for service={service_name} with query='{search_query}'")
    # LogQL syntax: {service="app"} |= "search_query"
    query = f'{{service="{service_name}"}} |= `{search_query}`'
    return query_loki(query, start_time, end_time, limit)

if __name__ == "__main__":
    mcp.run()
