-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL
);

-- +goose Down
DROP TABLE clients;
DROP EXTENSION IF EXISTS "uuid-ossp"