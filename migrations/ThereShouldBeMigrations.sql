CREATE TABLE IF NOT EXISTS posts(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    updated datetime NOT NULL,
    description text NOT NULL,
    image varbinary(max)
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    role text NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    description text,
    image varbinary(max)
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

СREATE TABLE IF NOT EXISTS comments(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    updated datetime NOT NULL,
    description text NOT NULL,
    image varbinary(max)
);
