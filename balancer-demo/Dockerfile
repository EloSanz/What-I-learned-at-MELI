FROM ghcr.io/astral-sh/uv:python3.12-alpine
WORKDIR /app

# Copiamos archivos de configuración de dependencias
COPY pyproject.toml uv.lock ./

# Sincronizamos las dependencias
RUN uv sync --frozen --no-dev

# Copiamos el código de la aplicación
COPY main.py .

# Añadimos el entorno virtual al PATH para ejecutar uvicorn directamente
ENV PATH="/app/.venv/bin:$PATH"

# Comando por defecto para correr FastAPI
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
