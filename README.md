# Delivery as a Service (DaaS) Backend - Arquitectura y Plan

Este repositorio contiene el motor de Delivery as a Service (DaaS) construido con Go. 
Ha sido diseñado siguiendo estrictamente los principios de **Mantenibilidad, Rendimiento, Seguridad y Resiliencia**.

## Convención de Trabajo (Branching Strategy)
Toda la implementación se realizará en **ramas secundarias** (`feature/*`), aplicando un riguroso *Code Review* enfocado en manejo de errores concurrentes, fugas de memoria y arquitectura limpia. Se requiere testing automático/manual antes de realizar el *merge* hacia la rama `main`.

## Arquitectura Limpia (Estructura)
El sistema emplea Inyección de Dependencias (DI) para separar el transporte HTTP, la lógica de negocio y el acceso a base de datos.

```text
/
├── cmd/
│   └── api/
│       └── main.go           # Entrypoint: Configuración DI, Handlers y Gin Web Server
├── internal/
│   ├── config/               # Carga de entorno segura (.env)
│   ├── models/               # Modelos de Datos GORM (Store, Driver, Order)
│   ├── repository/           # Capa de Datos (Interfaces e Implementación de Postgres)
│   ├── service/              # Capa de Negocio (Interfaces e Implementación)
│   └── handler/              # Capa de Transporte (Controladores Gin HTTP)
├── database/                 # Conexión principal de GORM (PostgreSQL)
├── Dockerfile                # Imagen multiestapa (builder + alpine base) minimizada
├── docker-compose.yml        # Orquestación app + bases de datos
└── .github/workflows/deploy.yml # Pipeline CI/CD para test -> build -> deploy (SSH VPS)
```

## Tecnologías Principales
*   **Lenguaje:** Go (Golang) 1.22+
*   **Web Framework:** Gin Gonic v1.9+
*   **ORM:** GORM (PostgreSQL Driver)
*   **Base de Datos:** PostgreSQL
*   **Infraestructura:** Docker & Docker Compose
*   **CI/CD:** GitHub Actions

## Detalles de Capas

### 1. Modelos (Entidades Base)
1.  `Store`: Entidad de comercio asociado a las órdenes. Posee coordenadas geográficas base.
2.  `Driver`: Repartidor con estado (Disponible, Ocupado) y rastreo de última ubicación.
3.  `Order`: Pedido central. Contiene transiciones de estado estrictas y referencias a su Store y Driver asignado.

### 2. Endpoints (API REST)
*   `POST /orders`: Emisión de un nuevo pedido por parte de un Store.
*   `GET /orders/nearby`: Retorna pedidos disponibles en un rango geográfico determinado (basado en Haversine o PostGIS subyacente).
*   `PATCH /orders/:id/status`: Transición autorizada del ciclo de vida del pedido.

### 3. Consideraciones de Performance y Seguridad
*   Middleware de control de pánicos para garantizar la disponibilidad del servicio.
*   Aseguramiento del cierre (`defer`) eficiente de descriptores de archivos, rows de SQL y cuerpos de respuesta HTTP.
*   Prevención de Goroutine Leaks al manejar contextos a lo largo de toda la petición GORM.
\n<!-- Sync test at Mon Mar 16 11:49:30 CET 2026 -->
