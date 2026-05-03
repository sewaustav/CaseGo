CREATE TABLE subscription_info (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    subscription INTEGER NOT NULL CHECK (subscription >= 0 AND subscription <=2),
    first_payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    count_of_renewal INTEGER NOT NULL DEFAULT 0,
    is_auto_renew BOOL DEFAULT FALSE,
    last_payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    canceled_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE payment_info (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    subscription_id BIGINT REFERENCES subscription_info(id) ON DELETE SET NULL,
    transaction_id VARCHAR(1000),
    price BIGINT NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    status VARCHAR(50),
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_subscription_user_id ON subscription_info (user_id);
CREATE INDEX idx_payment_user_id_date ON payment_info (user_id, payment_date DESC);
