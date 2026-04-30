# 🚀 Smart Delivery App (Grab-like)

A microservices-based smart delivery platform built with Flutter, Java Spring Boot, PostgreSQL, Redis, and Kafka — deployed on Kubernetes.

## Tech Stack

| Layer | Technology |
|---|---|
| Mobile | Flutter 3.x |
| API Gateway | Spring Cloud Gateway |
| Backend | Java 17 + Spring Boot 3.x (Microservices) |
| Primary DB | PostgreSQL 15 + PostGIS |
| Cache / Realtime | Redis 7 |
| Message Broker | Apache Kafka |
| Container | Docker + Docker Compose |
| Orchestration | Kubernetes (K8s) |
| Monitoring | Prometheus + Grafana |
| CI/CD | GitHub Actions |

## Services

| Service | Port | Description |
|---|---|---|
| api-gateway | 8080 | Single entry point, routing, rate limiting |
| auth-service | 8081 | JWT auth, login, register, refresh token |
| user-service | 8082 | Customer & driver profile management |
| order-service | 8083 | Order lifecycle (create → deliver) |
| tracking-service | 8084 | Real-time GPS tracking via WebSocket |
| payment-service | 8085 | Payment processing |
| notification-service | 8086 | Push notifications (FCM) |

## Project Structure

```
grab/
├── mobile/                     # Flutter App
├── backend/
│   ├── pom.xml                 # Parent Maven POM
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── order-service/
│   ├── tracking-service/
│   ├── payment-service/
│   └── notification-service/
├── database/
│   ├── schema.sql              # DDL
│   └── seed.sql                # Initial data
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml  # Local development
│   │   └── .env.example
│   └── k8s/                    # Kubernetes manifests
│       ├── namespaces/
│       ├── deployments/
│       ├── services/
│       ├── ingress/
│       ├── configmaps/
│       ├── secrets/
│       └── hpa/
├── monitoring/
│   ├── prometheus/
│   └── grafana/
└── .github/
    └── workflows/
        └── ci.yml
```

## Quick Start (Local)

### Prerequisites
- Docker Desktop
- Java 17+
- Flutter 3.x
- Maven 3.9+

### 1. Start Infrastructure
```bash
make infra-up
```

### 2. Run All Backend Services
```bash
make backend-up
```

### 3. Run Flutter App
```bash
make mobile-run
```

### Useful Commands
```bash
make infra-up        # Start PostgreSQL, Redis, Kafka
make infra-down      # Stop infrastructure
make backend-up      # Start all Spring Boot services
make backend-down    # Stop all Spring Boot services
make logs            # Tail all logs
make db-migrate      # Run database migrations
```

## Environment Variables

Copy `.env.example` to `.env` and fill in values:
```bash
cp infra/docker/.env.example infra/docker/.env
```

## API Documentation

After starting services, Swagger UI is available at:
- API Gateway: http://localhost:8080/swagger-ui.html
- Auth Service: http://localhost:8081/swagger-ui.html

## Architecture Overview

```
Flutter App
    ↕ HTTPS / WebSocket
API Gateway (:8080)
    ├── Auth Service (:8081)
    ├── User Service (:8082)
    ├── Order Service (:8083)    →  Kafka  →  Notification Service (:8086)
    ├── Tracking Service (:8084) ↔  Redis Pub/Sub
    └── Payment Service (:8085)
             ↓
        PostgreSQL + Redis
```
