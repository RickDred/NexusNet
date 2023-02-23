CREATE TABLE IF NOT EXISTS stories(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    content text NOT NULL,
    visible boolean NOT NULL
);
