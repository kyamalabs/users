-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2024-02-08T06:41:27.431Z

CREATE TABLE "profiles" (
  "wallet_address" varchar PRIMARY KEY,
  "gamer_tag" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "referrals" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "referrer" varchar,
  "referee" varchar,
  "referred_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "referrals" ADD FOREIGN KEY ("referrer") REFERENCES "profiles" ("wallet_address");

ALTER TABLE "profiles" ADD FOREIGN KEY ("wallet_address") REFERENCES "referrals" ("referee");
