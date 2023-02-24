CREATE TABLE IF NOT EXISTS direct(
                                      id bigserial PRIMARY KEY,
                                      user1 bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                      user2 bigint NOT NULL REFERENCES users ON DELETE CASCADE
);
