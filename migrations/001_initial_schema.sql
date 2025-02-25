-- +goose Up
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER NOT NULL REFERENCES brands(id),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    point_cost INTEGER NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    valid_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE redemptions (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id),
    total_points INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE redemption_items (
    id SERIAL PRIMARY KEY,
    redemption_id INTEGER NOT NULL REFERENCES redemptions(id),
    voucher_id INTEGER NOT NULL REFERENCES vouchers(id),
    quantity INTEGER NOT NULL,
    points_per_unit INTEGER NOT NULL,
    total_points INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add this INSERT statement after your CREATE TABLE statements
INSERT INTO customers (name, email, points)
VALUES ('Test Customer', 'test@example.com', 1000000);

-- +goose Down
DROP TABLE redemption_items;
DROP TABLE redemptions;
DROP TABLE customers;
DROP TABLE vouchers;
DROP TABLE brands;
