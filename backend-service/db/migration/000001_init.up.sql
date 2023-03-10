CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "hashad_password" varchar NOT NULL,
  "user_active" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "content" TEXT,
  "user_id" bigint
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");