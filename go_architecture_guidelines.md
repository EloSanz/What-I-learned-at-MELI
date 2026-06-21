# Lineamientos y Arquitectura de Microservicios en Go

Este documento detalla los principios de diseño, la estructura de carpetas y los estándares de desarrollo que siguen los servicios en Go (`items-service`, `orders-service` y futuros servicios) en esta plataforma de simulador de Mercado Libre.

---

## 1. Principios de Arquitectura

El diseño sigue una combinación de **Clean Architecture** (Arquitectura Limpia) y principios de **DDD (Domain-Driven Design)** simplificados para asegurar modularidad, desacoplamiento de dependencias y facilidad para escalar.

### Capas y Límites de Dependencia
El código se organiza en capas de adentro hacia afuera. Las dependencias solo pueden apuntar hacia adentro:
1.  **Dominio (Entidades y Reglas de Negocio):** Representado por las entidades (ej. `Item`, `Order`) y interfaces de repositorio. No depende de ningún framework ni librería externa (excepto las etiquetas de serialización y ORM directo).
2.  **Casos de Uso / Adaptadores (Repository & Handlers):** La lógica de cómo se recuperan los datos, cómo se manipulan transaccionalmente y cómo se procesan las peticiones HTTP (mediante controladores en Gin).
3.  **Frameworks y Drivers (Router & Database GORM):** La infraestructura externa. Base de datos PostgreSQL, enrutador HTTP Gin y librerías externas de red.

---

## 2. Estructura de Directorio Estándar

Cada microservicio sigue estrictamente esta estructura de archivos:

```text
├── cmd/
│   └── api/
│       └── main.go       # Punto de entrada y orquestador del servicio
├── internal/
│   ├── api/
│   │   └── router.go     # Configuración de Gin, middlewares y registro de rutas
│   ├── database/
│   │   └── db.go         # Conexión con GORM y configuración de base de datos
│   └── [domain]/         # Carpeta por cada dominio de negocio (ej: item, order)
│       ├── model.go      # Estructuras de dominio de GORM (Entidades)
│       ├── repo.go       # Acceso a base de datos (Interfaces y GORM repos)
│       └── handler.go    # Controladores HTTP (Input/Output HTTP, bindings)
└── pkg/
    └── web/
        └── response.go   # Utilidades HTTP genéricas y formatos estándar de JSON
```

### Descripción de Componentes:
*   **`cmd/api/main.go`**: Carga configuraciones de entorno, inicializa conexiones de base de datos, ejecuta migraciones automáticas (`AutoMigrate`), inicializa las capas del dominio (Repo -> Handler) y arranca el servidor HTTP escuchando señales de detención para un **Graceful Shutdown** (apagado ordenado).
*   **`internal/database/db.go`**: Establece la conexión del pool a Postgres. Centraliza valores por defecto y maneja errores de inicialización.
*   **`internal/[domain]/model.go`**: Declara la estructura del modelo SQL compatible con GORM y los formatos de serialización JSON.
*   **`internal/[domain]/repo.go`**: Implementa la interfaz del repositorio. Todas las operaciones de escritura que requieran atomicidad (como reducir stock) deben implementarse mediante transacciones nativas de GORM (`db.Transaction`) con bloqueos de fila (`FOR UPDATE`).
*   **`internal/[domain]/handler.go`**: Valida los datos entrantes del request utilizando el binding estructurado de Gin (`ShouldBindJSON`), delega la acción al repositorio y mapea la respuesta o el error usando las utilidades de `pkg/web`.
*   **`pkg/web/response.go`**: Mantiene un esquema JSON uniforme para todo el sistema:
    *   **Éxito:** `{"status": "success", "data": {...}, "message": "..."}`
    *   **Fallo:** `{"status": "error", "message": "..."}`

---

## 3. Flujo de Datos en una Solicitud

```mermaid
sequenceDiagram
    autonumber
    Client->>Router (Gin): GET /api/items/MLA43960787
    Router (Gin)->>Handler (HTTP): GetByID(c *gin.Context)
    Note over Handler: Valida parámetros y payloads
    Handler->>Repository: FindByID(id)
    Repository->>Database (GORM): Query Row
    Database (GORM)-->>Repository: Result/Error
    Repository-->>Handler: Entity/Error
    Note over Handler: Formatea respuesta estándar
    Handler-->>Client: JSON Response (200 OK / 404 Not Found)
```

---

## 4. Lineamientos de Desarrollo (Best Practices)

### A. Coherencia en Enrutamiento dentro de una VPC
Para facilitar el enrutamiento a nivel de Gateway (ALB, Nginx) mediante reglas basadas en paths, **todos los paths del servicio deben agruparse bajo el mismo prefijo del dominio**:
*   *Correcto:* `GET /api/items/:id` y `GET /api/items/health`
*   *Incorrecto:* `GET /items/:id` y `GET /health`

### B. Seguridad en Concurrencia y Transacciones
*   Cuando se realicen operaciones críticas sobre recursos compartidos (como el decremento de stock al comprar), se debe utilizar una transacción y adquirir un bloqueo de fila de base de datos (`FOR UPDATE`) para evitar condiciones de carrera (Race Conditions).

### C. Resiliencia
*   **Graceful Shutdown:** El servicio no debe morir abruptamente al recibir una señal de apagado (`SIGINT`, `SIGTERM`). Se deben tolerar de 5 a 10 segundos para dejar que las peticiones HTTP que ya están en vuelo terminen de procesarse.
*   **Tiempos de espera (Timeouts):** Al consumir APIs de otros microservicios (como de `orders-service` a `items-service`), se debe utilizar un cliente HTTP con un timeout explícito de máximo 5 segundos para evitar colgar hilos indefinidamente si el otro servicio está offline.
