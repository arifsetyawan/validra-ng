schema "public" {
  comment = "Validra Engine Schema"
}

table "resources" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "name" {
    type     = varchar(255)
    null     = false
  }
  column "description" {
    type     = text
    null     = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_resources_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "resource_actions" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "resource_id" {
    type = varchar(36)
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  foreign_key "fk_resource_actions_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.resources.column.id]
    on_delete   = CASCADE
  }
  
  index "idx_resource_actions_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "resource_abac_options" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "resource_id" {
    type = varchar(36)
    null = false
  }
  column "attribute" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  foreign_key "fk_resource_abac_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.resources.column.id]
    on_delete   = CASCADE
  }
  
  index "idx_resource_abac_options_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "resource_rebac_options" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "resource_id" {
    type = varchar(36)
    null = false
  }
  column "relation" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  foreign_key "fk_resource_rebac_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.resources.column.id]
    on_delete   = CASCADE
  }
  
  index "idx_resource_rebac_options_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "resource_sets" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_resource_sets_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "users" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "email" {
    type     = varchar(255)
    null     = false
  }
  column "entity" {
    type     = varchar(255)
    null     = false
  }
  column "attributes" {
    type = bytea
    null = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_users_email_unique" {
    columns = [column.email]
    unique = true
  }
  
  index "idx_users_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "user_sets" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "name" {
    type     = varchar(255)
    null     = false
  }
  column "criteria" {
    type = bytea
    null = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_user_sets_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "roles" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "name" {
    type     = varchar(255)
    null     = false
  }
  column "description" {
    type     = text
    null     = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_roles_deleted_at" {
    columns = [column.deleted_at]
  }
}

table "permissions" {
  schema = schema.public
  column "id" {
    type    = varchar(36)
    null    = false
  }
  column "role_id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = true
  }
  column "user_set_id" {
    type = varchar(36)
    null = true
  }
  column "resource_id" {
    type = varchar(36)
    null = true
  }
  column "resource_set_id" {
    type = varchar(36)
    null = true
  }
  column "effect" {
    type = varchar(10)
    null = false
  }
  column "conditions" {
    type = bytea
    null = true
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  
  primary_key {
    columns = [column.id]
  }
  
  foreign_key "fk_permissions_role" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_delete   = CASCADE
  }
  
  foreign_key "fk_permissions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = SET_NULL
  }
  
  foreign_key "fk_permissions_user_set" {
    columns     = [column.user_set_id]
    ref_columns = [table.user_sets.column.id]
    on_delete   = SET_NULL
  }
  
  foreign_key "fk_permissions_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.resources.column.id]
    on_delete   = SET_NULL
  }
  
  foreign_key "fk_permissions_resource_set" {
    columns     = [column.resource_set_id]
    ref_columns = [table.resource_sets.column.id]
    on_delete   = SET_NULL
  }
  
  index "idx_permissions_deleted_at" {
    columns = [column.deleted_at]
  }
} 