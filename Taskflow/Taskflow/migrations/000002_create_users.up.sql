CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_name_not_blank
        CHECK (length(trim(name)) > 0),

    CONSTRAINT users_email_not_blank
        CHECK (length(trim(email)) > 0)
);

CREATE UNIQUE INDEX users_email_unique
    ON users (lower(email));