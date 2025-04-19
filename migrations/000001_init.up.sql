CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE TASK_STATUS AS ENUM ('pending', 'done');

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    title VARCHAR(255) NOT NULL,
    status TASK_STATUS NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP
);
