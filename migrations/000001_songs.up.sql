CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS songs(
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name text not null,
    executor text not null,
    text jsonb not NULL,
    link text not null,
    release_date date not null,
    created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP
);


