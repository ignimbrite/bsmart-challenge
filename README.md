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
1) Copia variables de entorno:
   ```bash
   cp .env.example .env
   ```
2) Levanta PostgreSQL
   ```bash
   make compose-up
   ```
3) Ejecuta la API:
   ```bash
   make run
   ```

Health check: `http://localhost:8080/health`.# bsmart-challenge
