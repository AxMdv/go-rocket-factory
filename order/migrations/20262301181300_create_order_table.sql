-- +goose Up
CREATE TABLE orders (
    order_uuid uuid PRIMARY KEY,
    user_uuid text NOT NULL,
    part_uuids text[] NOT NULL,
    total_price_cents bigint NOT NULL CHECK (total_price_cents >= 0),
    transaction_uuid text,
    payment_method text,
    status text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT orders_part_uuids_not_empty CHECK (array_length(part_uuids, 1) > 0),
    CONSTRAINT orders_status_check CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED')),
    CONSTRAINT orders_payment_method_check CHECK (
        payment_method IS NULL OR payment_method IN ('CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY')
    ),
    CONSTRAINT orders_paid_requires_payment CHECK (
        status <> 'PAID' OR (transaction_uuid IS NOT NULL AND payment_method IS NOT NULL)
    )
);

-- +goose Down
DROP TABLE orders;
