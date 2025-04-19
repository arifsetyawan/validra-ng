# Validra Engine

A robust authorization and access control engine built with Go using Echo framework and SQLite.

## Project Structure

The project follows clean architecture principles and is structured as follows:

```
src/
  ├── config/           # Application configuration
  ├── internal/
  │   ├── delivery/     # API delivery layer
  │   │   └── http/     # HTTP handlers, DTOs, and middleware
  │   ├── domain/       # Domain models and repository interfaces
  │   ├── repository/   # Repository implementations
  │   └── service/      # Business logic services
  └── pkg/              # Shared packages
      ├── database/     # Database connection and migration
      ├── logger/       # Logging functionality
      └── validator/    # Request validation
```

## Getting Started

### Prerequisites

- Go 1.18 or later
- SQLite3

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/arifsetyawan/validra.git
   cd validra
   ```

2. Install the dependencies:
   ```
   go mod download
   ```

3. Build the application:
   ```
   go build -o validra ./src
   ```

### Configuration

The application can be configured using environment variables:

- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_READ_TIMEOUT`: Read timeout in seconds (default: 60)
- `SERVER_WRITE_TIMEOUT`: Write timeout in seconds (default: 60)
- `DB_PATH`: SQLite database file path (default: validra.db)

### Running the Application

```
./validra
```

Or directly with Go:

```
go run ./src
```

## API Endpoints

### Resources

- `POST /api/resources`: Create a new resource
- `GET /api/resources`: List resources
- `GET /api/resources/:id`: Get a specific resource
- `PUT /api/resources/:id`: Update a resource
- `DELETE /api/resources/:id`: Delete a resource

### Health Check

- `GET /health`: Check API health

## Development

### Adding a New Entity

To add a new entity to the system:

1. Define the domain model in `internal/domain/model.go`
2. Add repository interface in `internal/domain/repository.go`
3. Implement the repository in `internal/repository/`
4. Create a service in `internal/service/`
5. Define DTOs in `internal/delivery/http/dto/`
6. Create a handler in `internal/delivery/http/handler/`
7. Register routes in the main.go file

### Testing

Run the tests:
```
go test ./...
```

## License

This project is licensed under the MIT License.