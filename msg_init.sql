CREATE TABLE "messages" (
    "id" BIGINT PRIMARY KEY,
    "channel_id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "content" TEXT, 
    "media_id" TEXT,
    "message_type" VARCHAR(10) NOT NULL CHECK ("message_type" IN ('TEXT', 'MEDIA')),
    "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);

CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages(channel_id);
