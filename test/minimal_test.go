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

func TestMinimalGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-minimal-")
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

	// Verify that the expected files were generated
	for _, fileType := range []string{"server", "usecase", "repository"} {
		// Verify directory was created
		dirPath := filepath.Join(tmpDir, fileType)
		_, err = os.Stat(dirPath)
		require.NoError(t, err, "%s directory not created", fileType)

		// Verify service file was created
		serviceFilePath := filepath.Join(dirPath, "user_service_"+fileType+".go")
		_, err = os.Stat(serviceFilePath)
		require.NoError(t, err, "%s file not created", fileType)

		// Read file content
		content, err := os.ReadFile(serviceFilePath)
		require.NoError(t, err, "failed to read %s file", fileType)

		// Check basic content structure
		contentStr := string(content)
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
	}

	// Verify error codes file was created
	errorFilePath := filepath.Join(tmpDir, "error", "user_service_codes.go")
	_, err = os.Stat(errorFilePath)
	require.NoError(t, err, "error codes file not created")

	// Verify main init file was created
	mainFilePath := filepath.Join(tmpDir, "main_init.go")
	_, err = os.Stat(mainFilePath)
	require.NoError(t, err, "main init file not created")
}
