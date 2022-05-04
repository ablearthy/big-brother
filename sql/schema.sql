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