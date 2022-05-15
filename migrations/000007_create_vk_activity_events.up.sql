CREATE TYPE vk_platform AS ENUM ('mobile', 'iphone', 'ipad', 'android', 'wphone', 'windows', 'web');
CREATE TYPE vk_activity AS ENUM ('online', 'offline');

CREATE TABLE vk_activity_events (
    id SERIAL PRIMARY KEY,
    vk_owner_id INT NOT NULL,
    target_id INT NOT NULL,
    activity vk_activity NOT NULL,
    platform vk_platform,
    kicked_by_timeout BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);