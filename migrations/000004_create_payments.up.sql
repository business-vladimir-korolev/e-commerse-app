
CREATE TABLE e_commerce.payments (
    id                BIGSERIAL PRIMARY KEY,
    order_id          BIGINT NOT NULL REFERENCES e_commerce.orders(id),
    amount            INTEGER NOT NULL CHECK (amount >= 0),
    stripe_payment_id VARCHAR(255) NOT NULL UNIQUE,
    status            VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_payments_order ON e_commerce.payments (order_id);