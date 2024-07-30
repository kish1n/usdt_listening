-- +migrate Up

CREATE TABLE transfers (
   id UUID PRIMARY KEY,
   sender VARCHAR(42) NOT NULL,
   recipient VARCHAR(42) NOT NULL,
   value NUMERIC NOT NULL,
   timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE transfers;