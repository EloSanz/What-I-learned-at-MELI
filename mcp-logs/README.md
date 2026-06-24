# MCP Logs Server

This is a FastMCP server that connects to Loki to fetch and search logs from the `meli` Docker Compose stack.

## Connecting to IDEs

You can integrate this MCP server into AI assistants like **Cursor**, **Antigravity**, or **Claude Desktop**. 

Since this server is also running in Docker (`meli-mcp-logs` via SSE at `http://localhost:8005/sse`), some clients support HTTP/SSE connections. However, the most universally supported method for IDEs is running the server locally via `stdio`.

### Configuración Simplificada (Proxy HTTP)
Los IDEs actualmente solo soportan conexión vía `stdio` para los servidores MCP, por lo que no entienden los campos `type` o `url`.

Sin embargo, como el servidor ya corre en Docker, podés usar `uvx` para que `fastmcp` actúe como un proxy transparente que conecta el IDE (vía `stdio`) con el servidor remoto (vía `http`). 

Pegá este bloque en tu `mcp_config.json`:

```json
{
  "mcpServers": {
    "meli-logs-server": {
      "command": "uvx",
      "args": [
        "--from",
        "fastmcp-slim",
        "fastmcp",
        "run",
        "http://localhost:8005/sse"
      ]
    }
  }
}
```

¡Es la opción más limpia! El IDE usa comandos, pero **no tenés que definir paths locales, ni variables de entorno**. Toda la lógica y variables como `LOKI_URL` se manejan internamente dentro del contenedor Docker.
- **Type**: SSE
- **URL**: `http://localhost:8005/sse`

---

## Testing interactively
To test the server locally with a web UI:
```bash
cd mcp-logs
uv run fastmcp dev src/mcp_logs/server.py
```
