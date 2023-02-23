CREATE TABLE IF NOT EXISTS posts(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    title text NOT NULL,
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description text NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    role text NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    description text
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

CREATE TABLE IF NOT EXISTS comments(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    post_id bigint NOT NULL REFERENCES posts ON DELETE CASCADE,
    content text NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS stories(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    content text NOT NULL,
    visible boolean NOT NULL
);

