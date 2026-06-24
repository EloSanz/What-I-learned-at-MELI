FROM ghcr.io/astral-sh/uv:python3.12-alpine

WORKDIR /app

COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-dev

COPY *.py ./

ENV PATH="/app/.venv/bin:$PATH"

# Expose SSE transport port
EXPOSE 8000

CMD ["uv", "run", "--with", "fastmcp-slim[server]", "fastmcp", "run", "server.py", "--transport", "sse", "--host", "0.0.0.0", "--port", "8000"]
