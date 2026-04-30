-- ============================================================
-- Smart Delivery App — Database Schema
-- PostgreSQL 15 + PostGIS
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(100) NOT NULL,
    phone         VARCHAR(20)  UNIQUE NOT NULL,
    email         VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('CUSTOMER', 'DRIVER', 'ADMIN')),
    profile_image VARCHAR(500),
    fcm_token     VARCHAR(500),
    is_active     BOOLEAN      DEFAULT TRUE,
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

-- ============================================================
-- DRIVERS
-- ============================================================
CREATE TABLE drivers (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       UUID         UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vehicle_type  VARCHAR(20)  NOT NULL CHECK (vehicle_type IN ('MOTORCYCLE', 'CAR', 'VAN', 'TRUCK')),
    license_plate VARCHAR(20)  NOT NULL,
    vehicle_model VARCHAR(100),
    rating        DECIMAL(3,2) DEFAULT 5.00,
    total_trips   INTEGER      DEFAULT 0,
    is_available  BOOLEAN      DEFAULT FALSE,
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

-- ============================================================
-- DRIVER LOCATIONS (current position — also mirrored in Redis)
-- ============================================================
CREATE TABLE driver_locations (
    driver_id  UUID      PRIMARY KEY REFERENCES drivers(id) ON DELETE CASCADE,
    location   GEOMETRY(Point, 4326) NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_driver_locations_geom ON driver_locations USING GIST (location);

-- ============================================================
-- ORDERS
-- ============================================================
CREATE TABLE orders (
    id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id      UUID         NOT NULL REFERENCES users(id),
    driver_id        UUID         REFERENCES drivers(id),
    pickup_address   VARCHAR(500) NOT NULL,
    pickup_location  GEOMETRY(Point, 4326) NOT NULL,
    dropoff_address  VARCHAR(500) NOT NULL,
    dropoff_location GEOMETRY(Point, 4326) NOT NULL,
    status           VARCHAR(30)  NOT NULL DEFAULT 'PENDING'
                        CHECK (status IN ('PENDING','ACCEPTED','PICKUP','IN_TRANSIT','DELIVERED','CANCELLED')),
    total_price      DECIMAL(10,2),
    distance_km      DECIMAL(8,2),
    note             TEXT,
    cancelled_reason TEXT,
    created_at       TIMESTAMP    DEFAULT NOW(),
    updated_at       TIMESTAMP    DEFAULT NOW()
);

CREATE INDEX idx_orders_customer   ON orders(customer_id);
CREATE INDEX idx_orders_driver     ON orders(driver_id);
CREATE INDEX idx_orders_status     ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- ============================================================
-- ORDER ITEMS
-- ============================================================
CREATE TABLE order_items (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id    UUID         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    quantity    INTEGER      DEFAULT 1,
    weight_kg   DECIMAL(8,2),
    image_url   VARCHAR(500)
);

CREATE INDEX idx_order_items_order ON order_items(order_id);

-- ============================================================
-- PAYMENTS
-- ============================================================
CREATE TABLE payments (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id        UUID         UNIQUE NOT NULL REFERENCES orders(id),
    amount          DECIMAL(10,2) NOT NULL,
    method          VARCHAR(30)  NOT NULL CHECK (method IN ('CASH','CREDIT_CARD','QR_CODE','WALLET')),
    status          VARCHAR(20)  NOT NULL DEFAULT 'PENDING'
                        CHECK (status IN ('PENDING','COMPLETED','FAILED','REFUNDED')),
    transaction_ref VARCHAR(200),
    paid_at         TIMESTAMP,
    created_at      TIMESTAMP    DEFAULT NOW()
);

-- ============================================================
-- REVIEWS
-- ============================================================
CREATE TABLE reviews (
    id          UUID    PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id    UUID    UNIQUE NOT NULL REFERENCES orders(id),
    customer_id UUID    NOT NULL REFERENCES users(id),
    driver_id   UUID    NOT NULL REFERENCES drivers(id),
    rating      INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reviews_driver ON reviews(driver_id);

-- ============================================================
-- NOTIFICATIONS
-- ============================================================
CREATE TABLE notifications (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    body       TEXT         NOT NULL,
    type       VARCHAR(50)  NOT NULL,
    ref_id     UUID,
    is_read    BOOLEAN      DEFAULT FALSE,
    created_at TIMESTAMP    DEFAULT NOW()
);

CREATE INDEX idx_notifications_user    ON notifications(user_id);
CREATE INDEX idx_notifications_unread  ON notifications(user_id) WHERE is_read = FALSE;

-- ============================================================
-- REFRESH TOKENS
-- ============================================================
CREATE TABLE refresh_tokens (
    id         UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);

-- ============================================================
-- TRIGGER: updated_at auto-update
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at    BEFORE UPDATE ON users    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_drivers_updated_at  BEFORE UPDATE ON drivers  FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_orders_updated_at   BEFORE UPDATE ON orders   FOR EACH ROW EXECUTE FUNCTION update_updated_at();
