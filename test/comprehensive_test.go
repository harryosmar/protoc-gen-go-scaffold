package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComprehensiveGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-comprehensive-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test proto file with a more complex service definition
	complexProto := `syntax = "proto3";

package user;

option go_package = "github.com/harryosmar/protobuf-go/gen/user";

// UserService provides user management functionality
service UserService {
  // CreateUser creates a new user
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);

  // GetUser retrieves a user by ID
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // UpdateUser updates a user
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  
  // DeleteUser deletes a user
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  
  // ListUsers lists all users
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message CreateUserRequest {
  User user = 1;
}

message CreateUserResponse {
  User user = 1;
}

message GetUserRequest {
  uint32 id = 1;
}

message GetUserResponse {
  User user = 1;
}

message UpdateUserRequest {
  User user = 1;
}

message UpdateUserResponse {
  User user = 1;
}

message DeleteUserRequest {
  uint32 id = 1;
}

message DeleteUserResponse {
}

message ListUsersRequest {
  uint32 page = 1;
  uint32 limit = 2;
}

message ListUsersResponse {
  repeated User users = 1;
}

message User {
  uint32 id = 1;
  string name = 2;
  string email = 3;
}
`

	tmpProto := filepath.Join(tmpDir, "user.proto")
	err = os.WriteFile(tmpProto, []byte(complexProto), 0644)
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

	// Verify that the expected files were generated
	for _, fileType := range []string{"server", "usecase", "repository"} {
		// Verify directory was created
		dirPath := filepath.Join(tmpDir, fileType)
		_, err = os.Stat(dirPath)
		require.NoError(t, err, "%s directory not created", fileType)

		// Verify service file was created with correct naming convention
		serviceFilePath := filepath.Join(dirPath, "user_service_"+fileType+".go")
		_, err = os.Stat(serviceFilePath)
		require.NoError(t, err, "%s file not created", fileType)

		// Read file content
		content, err := os.ReadFile(serviceFilePath)
		require.NoError(t, err, "failed to read %s file", fileType)
		contentStr := string(content)

		// Check basic content structure
		assert.True(t, strings.Contains(contentStr, "package "+fileType),
			"%s file should have correct package", fileType)

		switch fileType {
		case "server":
			assert.True(t, strings.Contains(contentStr, "type UserServiceServer"),
				"%s file should define UserServiceServer type", fileType)
		case "usecase":
			assert.True(t, strings.Contains(contentStr, "type UserServiceUsecase interface"),
				"%s file should define UserServiceUsecase interface", fileType)
		case "repository":
			assert.True(t, strings.Contains(contentStr, "type ServiceRepository"),
				"%s file should define ServiceRepository interface", fileType)
		}

		// Check for all methods based on file type
		switch fileType {
		case "server", "usecase":
			assert.True(t, strings.Contains(contentStr, "CreateUser"),
				"%s file should contain CreateUser method", fileType)
			assert.True(t, strings.Contains(contentStr, "GetUser"),
				"%s file should contain GetUser method", fileType)
			assert.True(t, strings.Contains(contentStr, "UpdateUser"),
				"%s file should contain UpdateUser method", fileType)
			assert.True(t, strings.Contains(contentStr, "DeleteUser"),
				"%s file should contain DeleteUser method", fileType)
			assert.True(t, strings.Contains(contentStr, "ListUsers"),
				"%s file should contain ListUsers method", fileType)
		case "repository":
			assert.True(t, strings.Contains(contentStr, "Create"),
				"%s file should contain Create method", fileType)
			assert.True(t, strings.Contains(contentStr, "GetById"),
				"%s file should contain GetById method", fileType)
			assert.True(t, strings.Contains(contentStr, "Update"),
				"%s file should contain Update method", fileType)
			assert.True(t, strings.Contains(contentStr, "Delete"),
				"%s file should contain Delete method", fileType)
			assert.True(t, strings.Contains(contentStr, "GetPerPage"),
				"%s file should contain GetPerPage method", fileType)
		}
	}

	// Verify error codes file was created
	errorFilePath := filepath.Join(tmpDir, "error", "user_service_codes.go")
	_, err = os.Stat(errorFilePath)
	require.NoError(t, err, "error codes file not created")

	// Read error codes file content
	errorContent, err := os.ReadFile(errorFilePath)
	require.NoError(t, err, "failed to read error codes file")
	errorContentStr := string(errorContent)

	// Check error codes content
	assert.True(t, strings.Contains(errorContentStr, "package error"),
		"error codes file should have correct package")
	assert.True(t, strings.Contains(errorContentStr, "func InitUser"),
		"error codes file should contain init function")
	assert.True(t, strings.Contains(errorContentStr, "ErrUserNotFound"),
		"error codes file should contain error constants")

	// Verify main init file was created
	mainFilePath := filepath.Join(tmpDir, "main_init.go")
	_, err = os.Stat(mainFilePath)
	require.NoError(t, err, "main init file not created")

	// Read main init file content
	mainContent, err := os.ReadFile(mainFilePath)
	require.NoError(t, err, "failed to read main init file")
	mainContentStr := string(mainContent)

	// Check main init content
	assert.True(t, strings.Contains(mainContentStr, "package main"),
		"main init file should have correct package")
	assert.True(t, strings.Contains(mainContentStr, "func init"),
		"main init file should contain init function")
	assert.True(t, strings.Contains(mainContentStr, "InitUser"),
		"main init file should call InitUser")
}
