-- +migrate Up

CREATE TABLE transfers (
   id SERIAL PRIMARY KEY,
   from_address VARCHAR(42) NOT NULL,
   to_address VARCHAR(42) NOT NULL,
   value NUMERIC NOT NULL,
   timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
