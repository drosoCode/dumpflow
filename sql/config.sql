DROP TABLE IF EXISTS "sites";
CREATE TABLE "public"."sites" (
    "db_name" text PRIMARY KEY,
    "link" text NOT NULL,
    "update_date" timestamp NOT NULL,
    "auto_update" boolean NOT NULL,
    "enabled" boolean NOT NULL
) WITH (oids = false);
