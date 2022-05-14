CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR (16) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    inviter_id serial NOT NULL
);

CREATE TABLE invite_codes (
    user_id serial PRIMARY KEY NOT NULL,
    invite_code varchar (10) UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_tokens (
    user_id SERIAL PRIMARY KEY,
    access_token VARCHAR (100),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE vk_tokens (
    access_token varchar (100) PRIMARY KEY,
    vk_user_id INTEGER NOT NULL
);

CREATE TABLE vk_messages (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users (id),
    message_id INTEGER NOT NULL,
    message JSONB
);

CREATE TYPE vk_message_event_type AS ENUM ('new', 'edit', 'delete');

CREATE TABLE vk_message_events (
    id SERIAL PRIMARY KEY,
    internal_message_id INTEGER NOT NULL REFERENCES vk_messages (id),
    m_type vk_message_event_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);