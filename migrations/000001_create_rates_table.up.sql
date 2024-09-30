-- Migration: Create rates table

CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    token_a VARCHAR(10) NOT NULL,
    token_b VARCHAR(10) NOT NULL,
    ask VARCHAR(100) NOT NULL,
    bid VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);