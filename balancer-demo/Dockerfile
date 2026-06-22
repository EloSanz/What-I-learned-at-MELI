FROM ghcr.io/astral-sh/uv:python3.12-alpine
WORKDIR /app

# Copy dependency config files
COPY pyproject.toml uv.lock ./

# Sync dependencies
RUN uv sync --frozen --no-dev

# Copy application code
COPY app/ ./app/

# Add virtual env to PATH
ENV PATH="/app/.venv/bin:$PATH"

# Default command to run FastAPI
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]

