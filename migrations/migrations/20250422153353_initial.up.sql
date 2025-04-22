-- set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'Validra Engine Schema';
-- create "resources" table
CREATE TABLE "public"."resources" ("id" character varying(36) NOT NULL, "name" character varying(255) NOT NULL, "description" text NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_resources_deleted_at" to table: "resources"
CREATE INDEX "idx_resources_deleted_at" ON "public"."resources" ("deleted_at");
-- create "resource_sets" table
CREATE TABLE "public"."resource_sets" ("id" character varying(36) NOT NULL, "name" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_resource_sets_deleted_at" to table: "resource_sets"
CREATE INDEX "idx_resource_sets_deleted_at" ON "public"."resource_sets" ("deleted_at");
-- create "roles" table
CREATE TABLE "public"."roles" ("id" character varying(36) NOT NULL, "name" character varying(255) NOT NULL, "description" text NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deleted_at");
-- create "users" table
CREATE TABLE "public"."users" ("id" character varying(36) NOT NULL, "email" character varying(255) NOT NULL, "entity" character varying(255) NOT NULL, "attributes" bytea NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- create index "idx_users_email_unique" to table: "users"
CREATE UNIQUE INDEX "idx_users_email_unique" ON "public"."users" ("email");
-- create "user_sets" table
CREATE TABLE "public"."user_sets" ("id" character varying(36) NOT NULL, "name" character varying(255) NOT NULL, "criteria" bytea NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_user_sets_deleted_at" to table: "user_sets"
CREATE INDEX "idx_user_sets_deleted_at" ON "public"."user_sets" ("deleted_at");
-- create "permissions" table
CREATE TABLE "public"."permissions" ("id" character varying(36) NOT NULL, "role_id" character varying(36) NOT NULL, "user_id" character varying(36) NULL, "user_set_id" character varying(36) NULL, "resource_id" character varying(36) NULL, "resource_set_id" character varying(36) NULL, "effect" character varying(10) NOT NULL, "conditions" bytea NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_permissions_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "fk_permissions_resource_set" FOREIGN KEY ("resource_set_id") REFERENCES "public"."resource_sets" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "fk_permissions_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_permissions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "fk_permissions_user_set" FOREIGN KEY ("user_set_id") REFERENCES "public"."user_sets" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- create index "idx_permissions_deleted_at" to table: "permissions"
CREATE INDEX "idx_permissions_deleted_at" ON "public"."permissions" ("deleted_at");
-- create "resource_abac_options" table
CREATE TABLE "public"."resource_abac_options" ("id" character varying(36) NOT NULL, "resource_id" character varying(36) NOT NULL, "attribute" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_resource_abac_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_resource_abac_options_deleted_at" to table: "resource_abac_options"
CREATE INDEX "idx_resource_abac_options_deleted_at" ON "public"."resource_abac_options" ("deleted_at");
-- create "resource_actions" table
CREATE TABLE "public"."resource_actions" ("id" character varying(36) NOT NULL, "resource_id" character varying(36) NOT NULL, "name" character varying(255) NOT NULL, "description" text NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_resource_actions_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_resource_actions_deleted_at" to table: "resource_actions"
CREATE INDEX "idx_resource_actions_deleted_at" ON "public"."resource_actions" ("deleted_at");
-- create "resource_rebac_options" table
CREATE TABLE "public"."resource_rebac_options" ("id" character varying(36) NOT NULL, "resource_id" character varying(36) NOT NULL, "relation" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_resource_rebac_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_resource_rebac_options_deleted_at" to table: "resource_rebac_options"
CREATE INDEX "idx_resource_rebac_options_deleted_at" ON "public"."resource_rebac_options" ("deleted_at");
