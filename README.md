# BSMART Challenge

API REST + WebSocket para gestión de productos y categorías en Go.

## Stack
- Go 1.22
- Gin (HTTP)
- GORM + PostgreSQL
- Gorilla/WebSocket

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

## Endpoints
- `GET /api/products` — Lista productos
- `GET /health` — Health check