CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "forum" CASCADE;
DROP TABLE IF EXISTS "forum" CASCADE;

CREATE UNLOGGED TABLE IF NOT EXISTS "user"
(
    "id"       BIGSERIAL NOT NULL PRIMARY KEY,
    "nickname" CITEXT    NOT NULL UNIQUE,
    "fullname" CITEXT    NOT NULL,
    "about"    TEXT,
    "email"    CITEXT    NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS "forum"
(
    "id"      BIGSERIAL NOT NULL PRIMARY KEY,
    "title"   TEXT      NOT NULL,
    "user"    CITEXT    NOT NULL,
    "slug"    CITEXT    NOT NULL UNIQUE,
    "posts"   BIGINT DEFAULT 0,
    "threads" INT    DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS "thread"
(
    "id"      BIGSERIAL NOT NULL PRIMARY KEY,
    "title"   TEXT      NOT NULL,
    "author"  CITEXT    NOT NULL,
    "forum"   CITEXT,
    "message" TEXT      NOT NULL,
    "votes"   INT         DEFAULT 0,
    "slug"    CITEXT,
    "created" TIMESTAMPTZ DEFAULT now()
);
