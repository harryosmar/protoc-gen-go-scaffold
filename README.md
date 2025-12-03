# protoc-gen-go-scaffold

A Protocol Buffers compiler plugin that generates clean architecture scaffolding for Go applications.

## Features

Automatically generates:
- **Handler Layer** - HTTP handlers with DTOs for API endpoints
- **Service Layer** - Business logic with validation
- **Repository Layer** - Database operations with CRUD methods
- **Pagination Support** - Built-in pagination for list operations

## Installation

```bash
go install github.com/harryosmar/protoc-gen-go-scaffold@latest
```

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

### 1. Define Your Proto File

Create a proto file with entity messages (suffix with `Entity`):

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/yourproject/gen/user";

// UserEntity will generate full CRUD scaffold
message UserEntity {
  int64 id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
}

// UserDTO for API requests/responses
message UserDTO {
  string name = 1;
  string email = 2;
}
```

### 2. Generate Code

```bash
protoc -I./proto \
  --go_out=./gen --go_opt=paths=source_relative \
  --go-scaffold_out=./gen --go-scaffold_opt=paths=source_relative \
  proto/user.proto
```

### 3. Generated Files

The plugin generates three files:

```
gen/
├── user.pb.go                    # Standard protobuf (from protoc-gen-go)
├── user_handler.scaffold.go      # HTTP handlers
├── user_service.scaffold.go      # Business logic
└── user_repository.scaffold.go   # Database operations
```

## Generated Code Structure

### Handler Layer

```go
type UserEntityHandler struct {
    service *UserEntityService
}

// CRUD HTTP handlers
func (h *UserEntityHandler) CreateUserEntity(w http.ResponseWriter, r *http.Request)
func (h *UserEntityHandler) GetUserEntity(w http.ResponseWriter, r *http.Request)
func (h *UserEntityHandler) UpdateUserEntity(w http.ResponseWriter, r *http.Request)
func (h *UserEntityHandler) DeleteUserEntity(w http.ResponseWriter, r *http.Request)
func (h *UserEntityHandler) ListUserEntitys(w http.ResponseWriter, r *http.Request)
```

### Service Layer

```go
type UserEntityService struct {
    repo *UserEntityRepository
}

// Business logic methods
func (s *UserEntityService) Create(ctx context.Context, dto *UserEntityDTO) (*UserEntity, error)
func (s *UserEntityService) GetByID(ctx context.Context, id int64) (*UserEntity, error)
func (s *UserEntityService) Update(ctx context.Context, id int64, dto *UserEntityDTO) (*UserEntity, error)
func (s *UserEntityService) Delete(ctx context.Context, id int64) error
func (s *UserEntityService) List(ctx context.Context, page, pageSize int) ([]*UserEntity, int64, error)
```

### Repository Layer

```go
type UserEntityRepository struct {
    db *sql.DB
}

// Database operations
func (r *UserEntityRepository) Create(ctx context.Context, entity *UserEntity) error
func (r *UserEntityRepository) FindByID(ctx context.Context, id int64) (*UserEntity, error)
func (r *UserEntityRepository) Update(ctx context.Context, entity *UserEntity) error
func (r *UserEntityRepository) Delete(ctx context.Context, id int64) error
func (r *UserEntityRepository) FindAll(ctx context.Context, offset, limit int) ([]*UserEntity, error)
func (r *UserEntityRepository) Count(ctx context.Context) (int64, error)
```

## Example Usage

### Wire Up the Layers

```go
package main

import (
    "database/sql"
    "net/http"
    
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    
    pb "github.com/yourproject/gen/user"
)

func main() {
    // Database connection
    db, _ := sql.Open("postgres", "your-connection-string")
    defer db.Close()
    
    // Initialize layers
    repo := pb.NewUserEntityRepository(db)
    service := pb.NewUserEntityService(repo)
    handler := pb.NewUserEntityHandler(service)
    
    // Setup routes
    r := mux.NewRouter()
    r.HandleFunc("/users", handler.CreateUserEntity).Methods("POST")
    r.HandleFunc("/users/{id}", handler.GetUserEntity).Methods("GET")
    r.HandleFunc("/users/{id}", handler.UpdateUserEntity).Methods("PUT")
    r.HandleFunc("/users/{id}", handler.DeleteUserEntity).Methods("DELETE")
    r.HandleFunc("/users", handler.ListUserEntitys).Methods("GET")
    
    http.ListenAndServe(":8080", r)
}
```

### API Endpoints

- `POST /users` - Create user
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user
- `GET /users?page=1&page_size=10` - List users with pagination

## Entity Detection

The plugin identifies entities by:
1. Message name ending with `Entity` (e.g., `UserEntity`, `ProductEntity`)
2. Custom proto options (future feature)

## Customization

The generated code includes TODO comments where you need to:
- Add field mappings between DTOs and entities
- Implement validation logic
- Complete SQL queries with actual column names

## Dependencies

Generated code requires:
- `github.com/gorilla/mux` - HTTP routing
- `database/sql` - Database operations

## Development

### Build the Plugin

```bash
cd protoc-gen-go-scaffold
go build -o protoc-gen-go-scaffold main.go
```

### Install Locally

```bash
go install
```

### Test

```bash
cd ../protobuf-go
make proto
```

## Roadmap

- [ ] Custom proto options for entity/DTO annotations
- [ ] Support for different databases (MySQL, SQLite, MongoDB)
- [ ] Generate database migrations
- [ ] Generate OpenAPI/Swagger documentation
- [ ] Support for gRPC handlers (in addition to HTTP)
- [ ] Custom validation rules
- [ ] Relationship handling (foreign keys, joins)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License
