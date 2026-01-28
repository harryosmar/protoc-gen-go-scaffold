package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a simplified test proto file without external dependencies
	simplifiedProto := `syntax = "proto3";

package user;

option go_package = "github.com/harryosmar/protobuf-go/gen/user";

// UserService provides user management functionality
service UserService {
  // CreateUser creates a new user
  rpc CreateUser(CreateUserRequestDTO) returns (CreateUserResponseDTO);

  // GetUser retrieves a user by ID
  rpc GetUser(GetUserRequestDTO) returns (GetUserResponse);

  // UpdateUser updates a user
  rpc UpdateUser(UpdateUserRequestDTO) returns (UpdateUserResponseDTO);

  // DeleteUser deletes a user
  rpc DeleteUser(DeleteUserRequestDTO) returns (DeleteUserResponseDTO);

  // ListUsers lists all users
  rpc ListUsers(ListUsersRequestDTO) returns (ListUsersResponseDTO);
}

message CreateUserRequestDTO {
  UserDTO user = 1;
}

message CreateUserResponseDTO {
  UserDTO user = 1;
}

message GetUserRequestDTO {
  uint32 id = 1;
}

message GetUserResponse {
  UserDTO user = 1;
}

message UpdateUserRequestDTO {
  UserDTO user = 1;
}

message UpdateUserResponseDTO {
  UserDTO user = 1;
}

message DeleteUserRequestDTO {
  uint32 id = 1;
}

message DeleteUserResponseDTO {}

message ListUsersRequestDTO {
  PaginationRequest pagination = 1;
}

message ListUsersResponseDTO {
  repeated UserDTO users = 1;
  PaginationResponse pagination = 2;
}

message UserDTO {
  uint32 id = 1;
  string name = 2;
  string email = 3;
}

message PaginationRequest {
  int32 page = 1;
  int32 limit = 2;
}

message PaginationResponse {
  int64 total = 1;
  int32 page = 2;
  int32 limit = 3;
}
`

	tmpProto := filepath.Join(tmpDir, "user.proto")
	err = os.WriteFile(tmpProto, []byte(simplifiedProto), 0644)
	require.NoError(t, err)

	// Build plugin
	pluginPath := filepath.Join(cwd, "protoc-gen-go-scaffold")
	buildCmd := exec.Command("go", "build", "-o", pluginPath, "../main.go")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)
	defer os.Remove(pluginPath)

	// Run protoc with our plugin
	cmd := exec.Command("protoc",
		"--plugin=protoc-gen-go-scaffold="+pluginPath,
		"--proto_path=",
		"--proto_path="+tmpDir,
		"--go-scaffold_out=base=github.com/harryosmar/protobuf-go,paths=source_relative:"+tmpDir,
		filepath.Base(tmpProto),
	)
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "protoc execution failed: %s", output)

	// Test standard layers (server, usecase, repository)
	fileTypes := []string{"server", "usecase", "repository"}
	for _, fileType := range fileTypes {
		// Verify service filename
		expectedFilename := fmt.Sprintf("user_service_%s.go", fileType)
		serviceFilePath := filepath.Join(tmpDir, fileType, expectedFilename)
		_, err = os.Stat(serviceFilePath)
		require.NoError(t, err, "service %s file %s not found", fileType, expectedFilename)

		// Read generated content
		generatedContent, err := os.ReadFile(serviceFilePath)
		require.NoError(t, err, "service %s file not found", fileType)

		// Check that expected file exists
		expectedPath := filepath.Join(cwd, "expected", fileType, fmt.Sprintf("user_service_%s.go", fileType))
		_, err = os.Stat(expectedPath)
		require.NoError(t, err, "expected file not found")

		// Normalize content by trimming spaces
		normalizedGenerated := strings.TrimSpace(string(generatedContent))

		// Check that the file exists and contains key expected content
		assert.Contains(t, normalizedGenerated, "package "+fileType, "generated %s file should have correct package", fileType)

		switch fileType {
		case "server":
			assert.Contains(t, normalizedGenerated, "type UserServiceServer", "generated %s file should contain UserServiceServer type", fileType)
			assert.Contains(t, normalizedGenerated, "func NewUserServiceServer", "generated %s file should contain NewUserServiceServer function", fileType)
		case "usecase":
			assert.Contains(t, normalizedGenerated, "type UserServiceUsecase interface", "generated %s file should contain UserServiceUsecase interface", fileType)
			assert.Contains(t, normalizedGenerated, "func NewUserServiceUsecase", "generated %s file should contain NewUserServiceUsecase function", fileType)
		case "repository":
			assert.Contains(t, normalizedGenerated, "type ServiceRepository", "generated %s file should contain ServiceRepository interface", fileType)
			assert.Contains(t, normalizedGenerated, "func NewUserServiceRepositoryMySQL", "generated %s file should contain NewUserServiceRepositoryMySQL function", fileType)
		}
	}

	// Test error codes file
	errorFilePath := filepath.Join(tmpDir, "error", "user_service_codes.go")
	_, err = os.Stat(errorFilePath)
	require.NoError(t, err, "error codes file not found")

	// Read generated error codes content
	generatedErrorContent, err := os.ReadFile(errorFilePath)
	require.NoError(t, err, "error codes file not found")

	// Check that the error codes file exists and contains key expected content
	normalizedGeneratedError := strings.TrimSpace(string(generatedErrorContent))
	assert.Contains(t, normalizedGeneratedError, "package error", "generated error codes file should have correct package")
	assert.Contains(t, normalizedGeneratedError, "func Init", "generated error codes file should contain init function")
	assert.Contains(t, normalizedGeneratedError, "ErrUserNotFound", "generated error codes file should contain error constants")

	// Test main init file
	mainFilePath := filepath.Join(tmpDir, "main_init.go")
	_, err = os.Stat(mainFilePath)
	require.NoError(t, err, "main init file not found")

	// Read generated main init content
	generatedMainContent, err := os.ReadFile(mainFilePath)
	require.NoError(t, err, "main init file not found")

	// Check that the main init file exists and contains key expected content
	normalizedGeneratedMain := strings.TrimSpace(string(generatedMainContent))
	assert.Contains(t, normalizedGeneratedMain, "package main", "generated main init file should have correct package")
	assert.Contains(t, normalizedGeneratedMain, "func init", "generated main init file should contain init function")
	assert.Contains(t, normalizedGeneratedMain, "InitUserServiceCode", "generated main init file should call InitUserServiceCode")
}
