from src.mcp_logs.server import get_logs

if __name__ == "__main__":
    print("Testing get_logs...")
    # Fetch logs from 2026-06-24T00:00:00Z to 2026-06-25T00:00:00Z
    # Gateway A logs
    result = get_logs("gateway-a", "2026-06-24T00:00:00Z", "2026-06-25T00:00:00Z", 5)
    print("Logs:")
    print(result)
