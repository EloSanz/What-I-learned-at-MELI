import os
import httpx
from datetime import datetime, timezone

LOKI_URL = os.getenv("LOKI_URL", "http://localhost:3100")

def _parse_timestamp(ts: str) -> str:
    """Convert an ISO 8601 timestamp string into Unix nanoseconds format for Loki."""
    try:
        # Assuming format like "2026-06-24T12:00:00Z"
        dt = datetime.fromisoformat(ts.replace("Z", "+00:00"))
        # Loki API expects nanoseconds
        return str(int(dt.timestamp() * 1_000_000_000))
    except ValueError:
        # If it fails, return as is, assuming user provided nanoseconds directly
        return ts

def query_loki(query: str, start_time: str, end_time: str, limit: int = 100) -> str:
    """
    Executes a LogQL query against Loki via the HTTP API.
    """
    url = f"{LOKI_URL}/loki/api/v1/query_range"
    
    start_ns = _parse_timestamp(start_time)
    end_ns = _parse_timestamp(end_time)
    
    params = {
        "query": query,
        "start": start_ns,
        "end": end_ns,
        "limit": limit,
        "direction": "forward"
    }
    
    try:
        with httpx.Client() as client:
            response = client.get(url, params=params, timeout=10.0)
            response.raise_for_status()
            
            data = response.json()
            results = data.get("data", {}).get("result", [])
            
            if not results:
                return "No logs found for the given criteria."
                
            formatted_logs = []
            for stream in results:
                labels = stream.get("stream", {})
                container = labels.get("service", "unknown")
                for timestamp, log_line in stream.get("values", []):
                    formatted_logs.append(f"[{container}] {log_line}")
                    
            return "\n".join(formatted_logs)
            
    except httpx.HTTPError as e:
        return f"Error contacting Loki: {str(e)}"
    except Exception as e:
        return f"Unexpected error: {str(e)}"
