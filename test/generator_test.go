package test

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerator tests file generation without attempting to compile the output
func TestGenerator(t *testing.T) {
	// Get absolute path to the test directory
	testDir, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Setup paths
	pluginDir := filepath.Dir(testDir)
	testdataDir := filepath.Join(testDir, "testdata")
	outputDir := filepath.Join(testDir, "output")
	protoFile := filepath.Join(testdataDir, "user.proto")

	// Create output directory
	os.RemoveAll(outputDir)
	defer func() {
	}()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Build the protoc-gen-go-scaffold plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-go-scaffold", "-mod=mod")
	cmd.Dir = pluginDir
	// Set environment variables to ensure Go version compatibility
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build plugin: %v, output: %s", err, output)
	}
	defer os.Remove(filepath.Join(pluginDir, "protoc-gen-go-scaffold"))

	// Add plugin directory to PATH
	path := os.Getenv("PATH")
	os.Setenv("PATH", pluginDir+string(filepath.ListSeparator)+path)
	defer os.Setenv("PATH", path)

	// Using the original user.proto file without modifications
	t.Logf("Using original user.proto file: %s", protoFile)

	// Include paths for protoc - include all necessary proto dependencies
	// Project root is one level up from test dir
	projectRoot := filepath.Dir(testDir)
	thirdPartyDir := filepath.Join(projectRoot, "third_party")

	// Find additional proto paths in Go modules
	protoIncludePaths := []string{
		testdataDir,   // First check test directory
		thirdPartyDir, // Local third_party directory
		filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", "github.com", "grpc-ecosystem", "grpc-gateway@v1.16.0", "third_party", "googleapis"), // For google/api protos
		filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", "github.com", "infobloxopen", "protoc-gen-gorm@v1.1.5", "third_party", "proto"),      // For more dependencies
		os.Getenv("GOPATH") + "/src", // Then GOPATH
	}

	t.Logf("Using proto import paths from local and Go modules")

	// Build include path arguments
	args := []string{}
	for _, path := range protoIncludePaths {
		args = append(args, "-I", path)
	}

	// Add output and proto file
	args = append(args,
		"--go-scaffold_out=templates_dir="+pluginDir+",base=github.com/example/myapp:"+outputDir,
		"--experimental_allow_proto3_optional",
		protoFile)

	// Run protoc with our plugin
	protocCmd := exec.Command("protoc", args...)

	if output, err := protocCmd.CombinedOutput(); err != nil {
		t.Fatalf("protoc failed: %v, output: %s", err, string(output))
	}

	// Verify directory structure and files
	verifyOutputDirs(t, outputDir)

	t.Log("All files generated successfully!")
}

// copyDir recursively copies a directory tree from src to dst
func copyDir(t *testing.T, src, dst string) {
	info, err := os.Stat(src)
	if err != nil {
		t.Logf("Warning: Failed to stat source directory: %v", err)
		return
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		t.Logf("Warning: Failed to create destination directory: %v", err)
		return
	}

	// Read directory entries
	items, err := ioutil.ReadDir(src)
	if err != nil {
		t.Logf("Warning: Failed to read source directory: %v", err)
		return
	}

	// Copy each item
	for _, item := range items {
		srcPath := filepath.Join(src, item.Name())
		dstPath := filepath.Join(dst, item.Name())

		if item.IsDir() {
			// Recursively copy directories
			copyDir(t, srcPath, dstPath)
		} else {
			// Copy files
			copyFile(t, srcPath, dstPath)
		}
	}
}

// copyFile copies a single file from src to dst
func copyFile(t *testing.T, src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		t.Logf("Warning: Failed to open source file: %v", err)
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		t.Logf("Warning: Failed to create destination file: %v", err)
		return
	}
	defer out.Close()

	// Copy the content
	if _, err := io.Copy(out, in); err != nil {
		t.Logf("Warning: Failed to copy file content: %v", err)
		return
	}

	// Set permissions to match source file
	if fi, err := os.Stat(src); err == nil {
		os.Chmod(dst, fi.Mode())
	}
}

func verifyOutputDirs(t *testing.T, outputDir string) {
	// Check required directories exist
	requiredDirs := []string{"repository", "usecase", "server", "error"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(outputDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", dirPath)
		} else {
			// Check if directory contains any files
			files, err := ioutil.ReadDir(dirPath)
			if err != nil {
				t.Errorf("Failed to read directory %s: %v", dirPath, err)
			} else if len(files) == 0 {
				t.Errorf("Directory %s is empty", dirPath)
			} else {
				t.Logf("Directory %s contains %d files", dir, len(files))
			}
		}
	}

	// Check if main.go was generated
	mainFile := filepath.Join(outputDir, "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Errorf("Expected main.go does not exist")
	} else {
		// Verify if file has content
		content, err := ioutil.ReadFile(mainFile)
		if err != nil {
			t.Errorf("Failed to read main.go: %v", err)
		} else if len(content) == 0 {
			t.Errorf("main.go is empty")
		} else {
			t.Logf("main.go has %d bytes", len(content))

			// Check if the base package parameter was used
			if !strings.Contains(string(content), "github.com/example/myapp") {
				t.Errorf("Base package parameter not applied in main.go")
			}
		}
	}
}
