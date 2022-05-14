CREATE TYPE vk_message_event_type AS ENUM ('new', 'edit', 'delete');

CREATE TABLE vk_message_events (
    id SERIAL PRIMARY KEY,
    internal_message_id INTEGER NOT NULL REFERENCES vk_messages (id),
    m_type vk_message_event_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);