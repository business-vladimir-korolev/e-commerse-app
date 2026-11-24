CREATE TABLE e_commerce.orders (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES e_commerce.users(id),
    status       VARCHAR(20) NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'paid', 'failed', 'cancelled')),
    total_amount INTEGER NOT NULL CHECK (total_amount >= 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_orders_user ON e_commerce.orders (user_id);

CREATE TABLE e_commerce.order_items (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT NOT NULL
               REFERENCES e_commerce.orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL
               REFERENCES e_commerce.products(id) ON DELETE RESTRICT,
    quantity   INTEGER NOT NULL CHECK (quantity > 0),
    unit_price INTEGER NOT NULL CHECK (unit_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_order_items_order ON e_commerce.order_items (order_id);
