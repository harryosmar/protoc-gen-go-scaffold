# Quick Start Guide - protoc-gen-go-scaffold

## Installation Steps

### 1. Build and Install the Plugin

```bash
cd /Users/harry/go/src/github.com/harryosmar/protoc-gen-go-scaffold

# Build the plugin
go build -o protoc-gen-go-scaffold main.go

# Install to GOPATH/bin
go install
```

### 2. Verify Installation

```bash
which protoc-gen-go-scaffold
# Should output: /Users/harry/go/bin/protoc-gen-go-scaffold (or similar)
```

## Testing the Plugin

### 1. Create a Test Proto File

Create `test/user.proto`:

```protobuf
syntax = "proto3";

package test;

option go_package = "github.com/harryosmar/protoc-gen-go-scaffold/test/gen";

// UserEntity - will generate full CRUD scaffold
message UserEntity {
  int64 id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
}

// UserDTO - for API requests/responses
message UserDTO {
  string name = 1;
  string email = 2;
}
```

### 2. Generate Code

```bash
mkdir -p test/gen

protoc -I./test \
  --go_out=./test/gen --go_opt=paths=source_relative \
  --go-scaffold_out=./test/gen --go-scaffold_opt=paths=source_relative \
  test/user.proto
```

### 3. Check Generated Files

```bash
ls -la test/gen/
```

You should see:
- `user.pb.go` - Standard protobuf messages
- `user_handler.scaffold.go` - HTTP handlers
- `user_service.scaffold.go` - Business logic
- `user_repository.scaffold.go` - Database operations

## Using in Your Project

### 1. Update Makefile

Add to your project's Makefile:

```makefile
.PHONY: scaffold

scaffold:
	@mkdir -p gen
	protoc -I./proto \
		--go_out=./gen --go_opt=paths=source_relative \
		--go-scaffold_out=./gen --go-scaffold_opt=paths=source_relative \
		proto/*.proto
	@echo "✓ Scaffold generated successfully"
```

### 2. Generate Scaffold

```bash
make scaffold
```

## Troubleshooting

### Plugin Not Found

If you get "protoc-gen-go-scaffold: program not found":

```bash
# Check if it's in your PATH
echo $PATH | grep $(go env GOPATH)/bin

# If not, add to your shell profile (~/.zshrc or ~/.bashrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Reload shell
source ~/.zshrc
```

### Go Version Mismatch

If you see version mismatch errors:

```bash
# Check your Go version
go version

# Update go.mod to match
go mod edit -go=1.23

# Retry installation
go install
```

## Next Steps

1. Customize the generated code with your actual field mappings
2. Implement validation logic in the service layer
3. Complete SQL queries in the repository layer
4. Wire up the layers in your main.go
5. Add database migrations

## Example Integration

See the main README.md for a complete example of wiring up the generated layers.
