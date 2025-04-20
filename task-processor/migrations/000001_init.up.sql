CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP NOT NULL DEFAULT now ()
);
