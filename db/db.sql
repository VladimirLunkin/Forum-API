CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "forum" CASCADE;
DROP TABLE IF EXISTS "thread" CASCADE;
DROP TABLE IF EXISTS "post" CASCADE;
DROP TABLE IF EXISTS "vote" CASCADE;
DROP TABLE IF EXISTS "forum_user" CASCADE;

DROP FUNCTION IF EXISTS thread_vote();
DROP FUNCTION IF EXISTS create_post();
DROP FUNCTION IF EXISTS create_thread();

DROP TRIGGER IF EXISTS "vote_insert" ON "vote";
DROP TRIGGER IF EXISTS "create_post" ON "post";
DROP TRIGGER IF EXISTS "create_thread" ON "thread";

CREATE UNLOGGED TABLE IF NOT EXISTS "user"
(
    "id"       BIGSERIAL NOT NULL PRIMARY KEY,
    "nickname" CITEXT COLLATE "ucs_basic"   NOT NULL UNIQUE,
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

CREATE UNLOGGED TABLE IF NOT EXISTS "post"
(
    "id"       BIGSERIAL NOT NULL PRIMARY KEY,
    "parent"   BIGINT             DEFAULT 0,
    "author"   CITEXT    NOT NULL,
    "message"  TEXT      NOT NULL,
    "isEdited" BOOL               DEFAULT false,
    "forum"    CITEXT,
    "thread"   INT,
    "created"  TIMESTAMPTZ        DEFAULT now(),
    "path"     BIGINT[]  NOT NULL DEFAULT '{0}'
);

CREATE UNLOGGED TABLE IF NOT EXISTS "vote"
(
    "id"     BIGSERIAL                       NOT NULL PRIMARY KEY,
    "user"   BIGINT REFERENCES "user" (id)   NOT NULL,
    "thread" BIGINT REFERENCES "thread" (id) NOT NULL,
    "voice"  INT,
    CONSTRAINT checks UNIQUE ("user", "thread")
);

CREATE FUNCTION thread_vote() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "thread"
    SET "votes"=(votes + new.voice)
    WHERE "id" = new.thread;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "vote_insert"
    AFTER INSERT
    ON "vote"
    FOR EACH ROW
EXECUTE PROCEDURE thread_vote();

CREATE FUNCTION create_post() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "forum"
    SET "posts" = posts + 1
    WHERE "slug" = new.forum;
    new.path = (SELECT "path" FROM "post" WHERE "id" = new.parent LIMIT 1) || new.id;
    INSERT INTO "forum_user" ("user", "forum") VALUES ((SELECT "id" FROM "user" WHERE NEW.author = nickname), (SELECT "id" FROM "forum" WHERE NEW.forum = slug));
    return new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER "create_post"
    BEFORE INSERT
    ON "post"
    FOR EACH ROW
EXECUTE PROCEDURE create_post();

CREATE FUNCTION create_thread() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "forum"
    SET "threads" = threads + 1
    WHERE "slug" = new.forum;
    INSERT INTO "forum_user" ("user", "forum") VALUES ((SELECT "id" FROM "user" WHERE NEW.author = nickname), (SELECT "id" FROM "forum" WHERE NEW.forum = slug));
    return new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER "create_thread"
    BEFORE INSERT
    ON "thread"
    FOR EACH ROW
EXECUTE PROCEDURE create_thread();

CREATE UNLOGGED TABLE IF NOT EXISTS "forum_user"
(
    "id"    BIGSERIAL                      NOT NULL PRIMARY KEY,
    "user"  BIGINT REFERENCES "user" (id)  NOT NULL,
    "forum" BIGINT REFERENCES "forum" (id) NOT NULL
);

vacuum analyse;
