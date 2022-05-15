CREATE TABLE vk_messages (
    id SERIAL PRIMARY KEY,
    vk_owner_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    message JSONB
);