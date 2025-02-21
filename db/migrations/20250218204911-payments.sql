
-- +migrate Up
CREATE TYPE "status" AS ENUM ('pending', 'success', 'failed');

CREATE TABLE payments (
    "id" SERIAL PRIMARY KEY,
    "order_id" VARCHAR(100) NOT NULL,
    "user_id" INTEGER NOT NULL,
    "payment_method_id" INTEGER NOT NULL REFERENCES payment_methods(id),
    "status" status NOT NULL,
    "transaction_id" VARCHAR(255) NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS payments;
