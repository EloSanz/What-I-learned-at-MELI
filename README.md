# Meli Architecture Clone 🚀

Este proyecto es una simulación de la arquitectura orientada a microservicios inspirada en el flujo de compras de Mercado Libre. Está compuesto por un Frontend en Next.js, un API Gateway balanceado, varios microservicios en Go y un stack de observabilidad (Grafana + Loki).

## 🛠️ Cómo levantar la infraestructura

Para encender todos los motores, bases de datos y microservicios backend, simplemente usá Docker Compose en la raíz del proyecto:

```bash
docker-compose up -d
```

> **Nota:** La primera vez puede tardar un poco mientras descarga las imágenes de Postgres, Go, Python, Nginx y Grafana.

Para levantar la interfaz gráfica (Frontend):

```bash
cd my_meli
npm install
npm run dev
```

---

## 🗺️ Mapa de Puertos y Servicios

Una vez que todo esté corriendo, estos son los puertos de acceso para cada pieza de la arquitectura:

### 🖥️ Frontend & UI
- **Frontend (Tienda):** [http://localhost:3000](http://localhost:3000) (Next.js)
- **pgAdmin (Base de Datos Visual):** [http://localhost:5050](http://localhost:5050)
  - *Email:* `admin@admin.com`
  - *Password:* `postgrespassword`
- **Grafana (Métricas y Logs):** [http://localhost:8003](http://localhost:8003)

### 🔀 API Gateway & Balanceador
- **Nginx Balancer:** `http://localhost:8080`
  - `/gateway`: Endpoint principal para crear órdenes (rutea a FastAPI).
  - `/api/auth/login`: Endpoint público para obtener tu JWT.
- **Instancias FastAPI (Gateway Interno):** `Puerto 8000` (interno en Docker, escalado a 2 instancias: `gateway-a` y `gateway-b`).

### ⚙️ Microservicios Core (Backend)
Estos servicios corren en Go y no están expuestos directamente al público en la arquitectura ideal, pero podés acceder localmente para debug:
- **Items Service:** `http://localhost:8081` (Gestiona stock y catálogos. Conectado a `items_db`).
- **Orders Service:** `http://localhost:8082` (Gestiona la compra. Conectado a `orders_db`).
- **Auth Service:** `http://localhost:8083` (Genera y valida JWT. Validado vía Nginx `auth_request`).

### 🗄️ Bases de Datos e Infraestructura
- **PostgreSQL Engine:** `localhost:5432`
  - *User:* `postgres`
  - *Password:* `postgrespassword`
  - *Bases lógicas:* `items_db` y `orders_db`
- **RabbitMQ Dashboard:** [http://localhost:15672](http://localhost:15672)
  - *User:* `user`
  - *Password:* `password`
  - *Nota:* Acá podés ver en tiempo real cómo viajan los mensajes asincrónicos entre `orders-service` y `items-service`.
- **Loki & Promtail (Logging):** Corriendo en background recolectando logs de todos los contenedores para mostrarlos en Grafana.

---

## 🛒 Flujo de Compra (Cómo probarlo)

1. Abrí el Frontend en `http://localhost:3000`.
2. Hacé clic en **"Comprar ahora"** y luego confirmá.
3. El frontend internamente:
   - Hace un `POST` a `/api/auth/login` para obtener un Token JWT.
   - Envía el *Payload* de compra al `/gateway` (Nginx en puerto 8080) usando el header `Authorization: Bearer <token>`.
4. Nginx valida el token contra `Auth Service` y rutea la petición a una de las instancias de FastAPI.
5. FastAPI orquesta la compra: verifica stock con `Items Service` y crea la orden con `Orders Service`.
6. La interfaz te mostrará la pantalla verde de éxito indicando qué instancia resolvió tu pedido.
