package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCaptureOutput generates files and saves them to the expected directory
func TestCaptureOutput(t *testing.T) {

	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-capture-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a simplified test proto file without external dependencies
	simplifiedProto := `syntax = "proto3";

package user;

option go_package = "github.com/harryosmar/protobuf-go/gen/user";

// UserService provides user management functionality
service UserService {
  // CreateUser creates a new user
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);

  // GetUser retrieves a user by ID
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message CreateUserResponse {
  uint32 id = 1;
  string name = 2;
  string email = 3;
}

message GetUserRequest {
  uint32 id = 1;
}

message GetUserResponse {
  uint32 id = 1;
  string name = 2;
  string email = 3;
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

	// Copy generated files to expected directory
	expectedDir := filepath.Join(cwd, "expected")

	// Create directories if they don't exist
	for _, dir := range []string{"server", "usecase", "repository", "error"} {
		os.MkdirAll(filepath.Join(expectedDir, dir), 0755)
	}

	// Copy files
	for _, fileType := range []string{"server", "usecase", "repository"} {
		generatedPath := filepath.Join(tmpDir, fileType, fmt.Sprintf("user_service_%s.go", fileType))
		expectedPath := filepath.Join(expectedDir, fileType, fmt.Sprintf("user_service_%s.go", fileType))

		content, err := os.ReadFile(generatedPath)
		require.NoError(t, err, "Failed to read generated %s file", fileType)

		err = os.WriteFile(expectedPath, content, 0644)
		require.NoError(t, err, "Failed to write expected %s file", fileType)
	}

	// Copy error codes file
	errorPath := filepath.Join(tmpDir, "error", "user_service_codes.go")
	if _, err := os.Stat(errorPath); err == nil {
		content, err := os.ReadFile(errorPath)
		require.NoError(t, err, "Failed to read generated error file")

		expectedErrorPath := filepath.Join(expectedDir, "error", "user_service_codes.go")
		err = os.WriteFile(expectedErrorPath, content, 0644)
		require.NoError(t, err, "Failed to write expected error file")
	}

	// Copy main init file
	mainPath := filepath.Join(tmpDir, "main_init.go")
	if _, err := os.Stat(mainPath); err == nil {
		content, err := os.ReadFile(mainPath)
		require.NoError(t, err, "Failed to read generated main file")

		expectedMainPath := filepath.Join(expectedDir, "main_init.go")
		err = os.WriteFile(expectedMainPath, content, 0644)
		require.NoError(t, err, "Failed to write expected main file")
	}
}
