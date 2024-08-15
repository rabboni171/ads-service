CREATE TABLE "apps" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL UNIQUE,
    "secret" text NOT NULL UNIQUE
);
