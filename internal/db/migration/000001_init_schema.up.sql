-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2024-01-24T02:27:50.809Z

CREATE TABLE "profiles" (
  "wallet_address" varchar PRIMARY KEY,
  "gamer_tag" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
