CREATE TABLE IF NOT EXISTS messages(
                                     id bigserial PRIMARY KEY,
                                     direct_id bigint NOT NULL REFERENCES direct ON DELETE CASCADE,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                     content text NOT NULL,
                                     sender_id bigint NOT NULL REFERENCES users ON DELETE CASCADE
);
