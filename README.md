# protoc-gen-go-scaffold

A Protocol Buffers compiler plugin that automatically generates a clean architecture scaffold (service → usecase → repository) for Go gRPC services.

## Features

- **Automatic Layer Generation**: Creates service, usecase, and repository packages
- **Clean Architecture**: Enforces separation of concerns
- **Dependency Injection**: Sets up proper dependency injection patterns
- **Customizable Templates**: Easy to modify for specific needs
- **Convention-Based**: Follows standard Go naming conventions

## Build

```bash
go build -o $GOPATH/bin/protoc-gen-go-scaffold
```

## Installation

```bash
go install github.com/harryosmar/protobuf-go/protoc-gen-go-scaffold@latest
```

## Usage

Add to your protoc command:

```bash
protoc --go-scaffold_out=. --go-scaffold_opt=paths=source_relative your_service.proto
```

### Example

For a proto service definition:
```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

The plugin will generate:
```
├── repository
│   └── user_repository.go
├── service
│   └── user_service.go
└── usecase
    └── user_usecase.go
```

## Generated Structure

### 1. Service Layer (`service/`)
- Implements the gRPC service interface
- Handles request validation and transformation
- Delegates business logic to usecase

### 2. Usecase Layer (`usecase/`)
- Contains business logic
- Orchestrates data access through repository
- Handles domain-specific validation

### 3. Repository Layer (`repository/`)
- Implements data access logic
- Interacts with databases, APIs, or other data sources
- Contains persistence-related code

## Customization

Modify the template files in the `templates/` directory to customize:
- Package structure
- Dependency injection patterns
- Additional interfaces
- Custom business logic stubs

## Templates

The plugin uses Go templates for generation. You can find them in:
- `templates/service.tmpl`
- `templates/usecase.tmpl`
- `templates/repository.tmpl`

### Example Customization

To add logging to all service methods:

```gotemplate
// service.tmpl
func (s *{{$.Service.GoName}}Service) {{.GoName}}(ctx context.Context, req *{{.Input.GoIdent}}) (*{{.Output.GoIdent}}, error) {
    log.Printf("Handling {{.GoName}} request")
    return s.usecase.{{.GoName}}(ctx, req)
}
```

## Requirements

- Go 1.20+
- Protocol Buffers compiler (protoc)
- Google's Go protobuf plugin (`protoc-gen-go`)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.
