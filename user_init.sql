CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "username" VARCHAR(50) NOT NULL,
    "email" VARCHAR(100) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
    "updated_at" TIMESTAMP NOT NULL DEFAULT current_timestamp
);


CREATE TABLE "channels" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL UNIQUE, 
    "type" VARCHAR(10) NOT NULL CHECK ("type" IN ('GROUP', 'DM')), 
    "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE "memberships" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "channel_id" BIGINT NOT NULL REFERENCES "channels" ("id") ON DELETE CASCADE,
    "online_until" TIMESTAMP, 
    "joined_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
    UNIQUE ("user_id", "channel_id") 
);

CREATE INDEX IF NOT EXISTS idx_memberships_user_channel ON memberships(user_id, channel_id);

CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships(user_id);

CREATE INDEX IF NOT EXISTS idx_memberships_channel_id ON memberships(channel_id);
