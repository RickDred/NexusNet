CREATE TABLE IF NOT EXISTS comments(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    post_id bigint NOT NULL REFERENCES posts ON DELETE CASCADE,
    content text NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL
);
