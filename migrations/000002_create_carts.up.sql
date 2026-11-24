CREATE TABLE e_commerce.carts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL UNIQUE
               REFERENCES e_commerce.users(id) ON DELETE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE e_commerce.cart_items (
    id         BIGSERIAL PRIMARY KEY,
    cart_id    BIGINT NOT NULL
               REFERENCES e_commerce.carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL
               REFERENCES e_commerce.products(id) ON DELETE CASCADE,
    quantity   INTEGER NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_id)
);
CREATE INDEX idx_cart_items_product ON e_commerce.cart_items (product_id);