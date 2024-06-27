CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL
);

CREATE TABLE "rooms" (
    "room_id" bigserial PRIMARY KEY,
    "room_name" varchar NOT NULL,
    "created_by" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "room_members" (
    "user_id" bigint,
    "room_id" bigint,
    PRIMARY KEY ("user_id", "room_id")
);

CREATE TABLE "Messages" (
    "message_id" bigserial PRIMARY KEY,
    "room_id" bigint,
    "user_id" bigint,
    "message_content" text,
    "timestamp" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX ON "users" ("id");

CREATE INDEX ON "rooms" ("room_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "room_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "room_members" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "Messages" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "Messages" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");