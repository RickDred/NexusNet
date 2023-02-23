CREATE TABLE IF NOT EXISTS posts(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    title text NOT NULL,
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description text NOT NULL
);
