DROP TABLE IF EXISTS "sites";
CREATE TABLE "public"."sites" (
    "id" BIGSERIAL PRIMARY KEY,
    "db_name" text NOT NULL,
    "link" text NOT NULL,
    "update_date" timestamp NOT NULL,
    "auto_update" boolean NOT NULL,
    "enabled" boolean NOT NULL
) WITH (oids = false);
