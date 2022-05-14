CREATE TABLE vk_messages (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users (id),
    message_id INTEGER NOT NULL,
    message JSONB
);