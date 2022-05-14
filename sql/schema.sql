-- Adminer 4.7.8 PostgreSQL dump

DROP TABLE IF EXISTS "badges";
CREATE TABLE "public"."badges" (
    "id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "name" text NOT NULL,
    "date" timestamp NOT NULL,
    "class" integer NOT NULL,
    "tag_based" boolean NOT NULL
) WITH (oids = false);


DROP TABLE IF EXISTS "comments";
CREATE TABLE "public"."comments" (
    "id" bigint NOT NULL,
    "post_id" bigint NOT NULL,
    "score" integer NOT NULL,
    "text" text NOT NULL,
    "creation_date" timestamp NOT NULL,
    "user_id" bigint NOT NULL,
    "content_license" text NOT NULL,
    "text_index" tsvector GENERATED ALWAYS AS (to_tsvector('english', "text")) STORED
) WITH (oids = false);


DROP TABLE IF EXISTS "post_history";
CREATE TABLE "public"."post_history" (
    "id" bigint NOT NULL,
    "post_history_type_id" integer NOT NULL,
    "post_id" bigint NOT NULL,
    "revision_guid" text NOT NULL,
    "creation_date" timestamp NOT NULL,
    "user_id" bigint NOT NULL,
    "comment" text NOT NULL,
    "text" text NOT NULL,
    "content_license" text NOT NULL,
    "text_index" tsvector GENERATED ALWAYS AS (to_tsvector('english', "text")) STORED
) WITH (oids = false);


DROP TABLE IF EXISTS "post_links";
CREATE TABLE "public"."post_links" (
    "id" bigint NOT NULL,
    "creation_date" timestamp NOT NULL,
    "post_id" bigint NOT NULL,
    "related_post_id" bigint NOT NULL,
    "link_type_id" integer NOT NULL
) WITH (oids = false);


DROP TABLE IF EXISTS "posts";
CREATE TABLE "public"."posts" (
    "id" bigint NOT NULL,
    "post_type_id" integer NOT NULL,
    "parent_id" bigint NOT NULL,
    "accepted_answer_id" integer NOT NULL,
    "creation_date" timestamp NOT NULL,
    "closed_date" timestamp NOT NULL,
    "score" integer NOT NULL,
    "view_count" integer NOT NULL,
    "body" text NOT NULL,
    "tags" text NOT NULL,
    "answer_count" integer NOT NULL,
    "comment_count" integer NOT NULL,
    "favorite_count" integer NOT NULL,
    "content_license" text NOT NULL,
    "body_index" tsvector GENERATED ALWAYS AS (to_tsvector('english', "body")) STORED
) WITH (oids = false);


DROP TABLE IF EXISTS "tags";
CREATE TABLE "public"."tags" (
    "id" bigint NOT NULL,
    "tag_name" text NOT NULL,
    "count" integer NOT NULL,
    "is_required" boolean NOT NULL,
    "is_moderator_only" boolean NOT NULL,
    "wiki_post_id" integer NOT NULL,
    "excerpt_post_id" integer NOT NULL
) WITH (oids = false);


DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "id" bigint NOT NULL,
    "reputation" int NOT NULL,
    "creation_date" timestamp NOT NULL,
    "display_name" text NOT NULL,
    "last_access_date" timestamp NOT NULL,
    "location" text NOT NULL,
    "website_url" text NOT NULL,
    "about_me" text NOT NULL,
    "views" int NOT NULL,
    "upvotes" int NOT NULL,
    "downvotes" int NOT NULL,
    "account_id" bigint NOT NULL,
    "profile_image_url" text NOT NULL
) WITH (oids = false);


DROP TABLE IF EXISTS "votes";
CREATE TABLE "public"."votes" (
    "id" bigint NOT NULL,
    "post_id" bigint NOT NULL,
    "vote_type_id" int NOT NULL,
    "creation_date" timestamp NOT NULL
) WITH (oids = false);
