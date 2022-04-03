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