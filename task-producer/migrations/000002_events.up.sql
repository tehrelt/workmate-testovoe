create table if not exists events (
    id uuid primary key,
    created_at timestamp not null default now ()
);
