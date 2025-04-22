env "dev" {
  // Dev environment configuration
  src = "migrations/schemas/schema.hcl"
  url = "postgres://postgres:postgres@localhost:5432/validra?sslmode=disable"
  dev = "postgres://postgres:postgres@localhost:5432/validra?sslmode=disable"
  migration {
    dir = "file://migrations/migrations"
    format = golang-migrate
  }
}

env "prod" {
  // Production environment configuration 
  src = "migrations/schemas/schema.hcl"
  url = "postgres://${env.DB_USER}:${env.DB_PASSWORD}@${env.DB_HOST}:${env.DB_PORT}/${env.DB_NAME}?sslmode=${env.DB_SSL_MODE}"
  migration {
    dir = "file://migrations/migrations"
    format = golang-migrate
  }
} 