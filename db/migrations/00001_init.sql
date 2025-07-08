-- +goose Up
-- +goose StatementBegin
CREATE TABLE contracts (
    id UUID PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_email TEXT NOT NULL CHECK (contact_email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
    monthly_limit NUMERIC(15,2) NOT NULL CHECK (monthly_limit >= 1000),
    current_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL CHECK (end_date > start_date),
    status VARCHAR(10) NOT NULL DEFAULT 'active' 
        CHECK (status IN ('active', 'expired', 'blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contracts;
-- +goose StatementEnd
