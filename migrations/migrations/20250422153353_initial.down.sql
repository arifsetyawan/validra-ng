-- reverse: create index "idx_resource_rebac_options_deleted_at" to table: "resource_rebac_options"
DROP INDEX "public"."idx_resource_rebac_options_deleted_at";
-- reverse: create "resource_rebac_options" table
DROP TABLE "public"."resource_rebac_options";
-- reverse: create index "idx_resource_actions_deleted_at" to table: "resource_actions"
DROP INDEX "public"."idx_resource_actions_deleted_at";
-- reverse: create "resource_actions" table
DROP TABLE "public"."resource_actions";
-- reverse: create index "idx_resource_abac_options_deleted_at" to table: "resource_abac_options"
DROP INDEX "public"."idx_resource_abac_options_deleted_at";
-- reverse: create "resource_abac_options" table
DROP TABLE "public"."resource_abac_options";
-- reverse: create index "idx_permissions_deleted_at" to table: "permissions"
DROP INDEX "public"."idx_permissions_deleted_at";
-- reverse: create "permissions" table
DROP TABLE "public"."permissions";
-- reverse: create index "idx_user_sets_deleted_at" to table: "user_sets"
DROP INDEX "public"."idx_user_sets_deleted_at";
-- reverse: create "user_sets" table
DROP TABLE "public"."user_sets";
-- reverse: create index "idx_users_email_unique" to table: "users"
DROP INDEX "public"."idx_users_email_unique";
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX "public"."idx_users_deleted_at";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_roles_deleted_at" to table: "roles"
DROP INDEX "public"."idx_roles_deleted_at";
-- reverse: create "roles" table
DROP TABLE "public"."roles";
-- reverse: create index "idx_resource_sets_deleted_at" to table: "resource_sets"
DROP INDEX "public"."idx_resource_sets_deleted_at";
-- reverse: create "resource_sets" table
DROP TABLE "public"."resource_sets";
-- reverse: create index "idx_resources_deleted_at" to table: "resources"
DROP INDEX "public"."idx_resources_deleted_at";
-- reverse: create "resources" table
DROP TABLE "public"."resources";
-- reverse: set comment to schema: "public"
COMMENT ON SCHEMA "public" IS NULL;
