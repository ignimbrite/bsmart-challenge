# BSMART Challenge

API REST + WebSocket para gestión de productos y categorías en Go.

## Stack
- Go 1.22
- Gin (HTTP)
- GORM + PostgreSQL
- Gorilla/WebSocket
- JWT (github.com/golang-jwt/jwt/v5)

## Requisitos
- Go 1.22+
- Docker y Docker Compose

## Configuración rápida
1) Variables de entorno:
   ```bash
   cp .env.example .env
   ```

2) Base de datos:
   ```bash
   make compose-up
   ```

3) Ejecuta la API:
   ```bash
   make run
   ```

Health check: `http://localhost:8080/health`.

## Endpoints principales
- Auth: `POST /api/auth/login` (JWT). Usuario seed en dev: `admin@bsmart.test` / `admin123`.
- Products: `GET /api/products`, `GET /api/products/:id`, `POST /api/products`, `PUT /api/products/:id`, `DELETE /api/products/:id`, `GET /api/products/:id/history?start=YYYY-MM-DD&end=YYYY-MM-DD`
- Categories: `GET /api/categories`, `POST /api/categories`, `PUT /api/categories/:id`, `DELETE /api/categories/:id`
- Search: `GET /api/search?type=product|category&q=&page=&page_size=&sort=`
- Health: `GET /health`

Notas:
- Rutas de escritura (POST/PUT/DELETE) requieren JWT de un usuario con rol `admin`.
- Paginación: `page`, `page_size` (máx 100), `sort` soporta `price_asc|price_desc|name_asc|name_desc|newest|oldest` en productos y `name_asc|name_desc|newest|oldest` en categorías.
- Historial: guarda cambios de `price`/`stock`; se puede filtrar por rango de fechas.
- WebSocket: `GET /ws` (upgrade). Eventos `product.created|updated|deleted` y `category.created|updated|deleted`.

En entorno `APP_ENV=development` se insertan datos de ejemplo (categorías/productos) y usuario admin para login.
