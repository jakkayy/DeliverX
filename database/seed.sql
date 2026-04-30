-- ============================================================
-- Seed Data for Development
-- ============================================================

-- Admin user (password: admin1234)
INSERT INTO users (id, name, phone, email, password_hash, role) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Admin', '0800000000', 'admin@grab.local',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'ADMIN');

-- Customer users (password: password123)
INSERT INTO users (id, name, phone, email, password_hash, role) VALUES
    ('00000000-0000-0000-0000-000000000002', 'Somchai Jaidee', '0811111111', 'somchai@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'CUSTOMER'),
    ('00000000-0000-0000-0000-000000000003', 'Malee Sriwan', '0822222222', 'malee@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'CUSTOMER');

-- Driver users (password: password123)
INSERT INTO users (id, name, phone, email, password_hash, role) VALUES
    ('00000000-0000-0000-0000-000000000004', 'Prasert Nakorn', '0833333333', 'prasert@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'DRIVER'),
    ('00000000-0000-0000-0000-000000000005', 'Wanchai Daeng', '0844444444', 'wanchai@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'DRIVER');

-- Drivers
INSERT INTO drivers (id, user_id, vehicle_type, license_plate, vehicle_model, rating, total_trips, is_available) VALUES
    ('00000000-0000-0000-0001-000000000001', '00000000-0000-0000-0000-000000000004',
     'MOTORCYCLE', 'กข-1234', 'Honda PCX', 4.85, 320, TRUE),
    ('00000000-0000-0000-0001-000000000002', '00000000-0000-0000-0000-000000000005',
     'CAR', 'คง-5678', 'Toyota Yaris', 4.72, 180, TRUE);

-- Driver locations (Bangkok area)
INSERT INTO driver_locations (driver_id, location) VALUES
    ('00000000-0000-0000-0001-000000000001', ST_SetSRID(ST_MakePoint(100.5018, 13.7563), 4326)),
    ('00000000-0000-0000-0001-000000000002', ST_SetSRID(ST_MakePoint(100.5230, 13.7400), 4326));
