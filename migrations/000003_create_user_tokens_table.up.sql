CREATE TABLE IF NOT EXISTS user_tokens (
    user_id SERIAL PRIMARY KEY,
    access_token VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES users(id)
);