# DeliverX

โปรเจคส่วนตัวที่ทำขึ้นมาเพื่อฝึก microservices architecture โดยจำลองแอปส่งของคล้ายๆ Grab / Lalamove ครับ

stack หลักเป็น Go + Gin สำหรับ backend, Flutter สำหรับ mobile และใช้ PostgreSQL, Redis, Kafka เป็น infrastructure — ทั้งหมดรันผ่าน Docker และ deploy บน Kubernetes

---

## Tech Stack

| Layer | Technology |
|---|---|
| Mobile | Flutter 3.x |
| API Gateway | Go + Gin |
| Backend | Go 1.22+ + Gin (Microservices) |
| Database | PostgreSQL 15 + PostGIS |
| Cache / Realtime | Redis 7 |
| Message Broker | Apache Kafka |
| Container | Docker + Docker Compose |
| Orchestration | Kubernetes |
| Monitoring | Prometheus + Grafana |
| CI/CD | GitHub Actions |

---

## Services

แบ่งออกเป็น 7 services ตาม domain ครับ

| Service | Port | หน้าที่ |
|---|---|---|
| api-gateway | 8080 | รับ request ทั้งหมด แล้ว proxy ไปแต่ละ service + validate JWT |
| auth-service | 8081 | register, login, refresh token, logout |
| user-service | 8082 | จัดการ profile ของลูกค้าและคนขับ |
| order-service | 8083 | สร้างและติดตาม order ตั้งแต่ต้นจนจบ |
| tracking-service | 8084 | รับ GPS จากคนขับแบบ real-time ผ่าน WebSocket + Redis GEO |
| payment-service | 8085 | จัดการการชำระเงิน |
| notification-service | 8086 | ฟัง event จาก Kafka แล้วยิง push notification ผ่าน FCM |

---

## Mobile Features

Flutter app มีหน้าหลักดังนี้ครับ

| หน้า | คำอธิบาย |
|---|---|
| Login / Register | ยืนยันตัวตนผ่าน auth-service |
| Home | หน้าหลักหลัง login |
| Create Order | สร้างคำสั่งซื้อใหม่ |
| Order Detail | ดูรายละเอียดและสถานะ order |
| Tracking | ติดตามคนขับแบบ real-time ผ่าน WebSocket |
| Profile | จัดการข้อมูลส่วนตัว |

State management: **Flutter BLoC** — Navigation: **go_router** — HTTP: **Dio + Retrofit**

---

## Project Structure

โครงสร้างของแต่ละ service จะเหมือนกันหมดครับ แบบนี้

```
<service>/
├── cmd/main.go          # entry point
├── config/              # อ่าน env vars
├── internal/
│   ├── handler/         # Gin routes
│   ├── service/         # business logic
│   ├── repository/      # query DB
│   └── model/           # struct + GORM
└── go.mod
```

ภาพรวมทั้งโปรเจค

```
DeliverX/
├── mobile/
├── backend/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── order-service/
│   ├── tracking-service/
│   ├── payment-service/
│   └── notification-service/
├── database/
│   ├── schema.sql
│   └── seed.sql
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   └── .env.example
│   └── k8s/
├── monitoring/
└── .github/workflows/
```

---

## รันบนเครื่องตัวเอง

ต้องมี Docker Desktop กับ Go 1.22+ ติดตั้งไว้ก่อนครับ Flutter ด้วยถ้าจะรัน mobile

**1. เปิด infrastructure (PostgreSQL, Redis, Kafka)**
```bash
make infra-up
```

**2. รัน backend ทั้งหมด**
```bash
make backend-up
```

**3. รัน Flutter app**
```bash
make mobile-run
```

ก่อนรันครั้งแรกอย่าลืม copy env file ด้วยนะครับ
```bash
cp infra/docker/.env.example infra/docker/.env
```

### คำสั่งที่ใช้บ่อย

```bash
make infra-up        # เปิด PostgreSQL, Redis, Kafka
make infra-down      # ปิด infrastructure
make backend-up      # เปิด services ทั้งหมด
make backend-down    # ปิด services ทั้งหมด
make logs            # ดู logs
make db-migrate      # migrate database
make tidy            # go mod tidy ทุก service พร้อมกัน
```

---

## Architecture

```
Flutter App
    ↕ HTTPS / WebSocket
API Gateway (:8080)
    ├── Auth Service (:8081)
    ├── User Service (:8082)
    ├── Order Service (:8083)    →  Kafka  →  Notification Service (:8086)
    ├── Tracking Service (:8084) ↔  Redis GEO / WebSocket
    └── Payment Service (:8085)  →  Kafka  →  Notification Service (:8086)
             ↓
        PostgreSQL + Redis
```
