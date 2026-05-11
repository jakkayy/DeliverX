# DeliverX — Smart Delivery App

แพลตฟอร์มรับส่งพัสดุแบบ real-time พัฒนาด้วยสถาปัตยกรรม Microservices ประกอบด้วย Flutter, Java Spring Boot, PostgreSQL, Redis และ Kafka — deploy บน Kubernetes

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

| Service | Port | คำอธิบาย |
|---|---|---|
| api-gateway | 8080 | จุดเข้าหลัก, routing, rate limiting |
| auth-service | 8081 | ยืนยันตัวตนด้วย JWT, login, register, refresh token |
| user-service | 8082 | จัดการโปรไฟล์ลูกค้าและคนส่งของ |
| order-service | 8083 | วงจรชีวิตคำสั่งซื้อ (สร้าง → จัดส่ง) |
| tracking-service | 8084 | ติดตาม GPS แบบ real-time ผ่าน WebSocket |
| payment-service | 8085 | ประมวลผลการชำระเงิน |
| notification-service | 8086 | ส่ง push notification ผ่าน FCM |

## Project Structure

```
DeliverX/
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
│   ├── schema.sql              # สร้างตารางฐานข้อมูล
│   └── seed.sql                # ข้อมูลเริ่มต้น
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml  # สำหรับรันบนเครื่อง local
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

### 1. เปิด Infrastructure
```bash
make infra-up
```

### 2. รัน Backend Services ทั้งหมด
```bash
make backend-up
```

### 3. รัน Flutter App
```bash
make mobile-run
```

### Useful Commands
```bash
make infra-up        # เปิด PostgreSQL, Redis, Kafka
make infra-down      # ปิด infrastructure
make backend-up      # เปิด Spring Boot services ทั้งหมด
make backend-down    # ปิด Spring Boot services ทั้งหมด
make logs            # ดู logs ทั้งหมด
make db-migrate      # รัน database migrations
```

## API Documentation

หลังจากเปิด services แล้ว เข้าดู Swagger UI ได้ที่:
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
