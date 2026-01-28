package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// TestTemplateRendering tests that the usecase template renders correctly with ORM data
func TestTemplateRendering(t *testing.T) {
	// Get the project root directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	// Path to the template file
	tmplPath := filepath.Join(cwd, "..", "templates", "usecase.tmpl")
	
	// Load the template file
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	// Define template functions
	funcMap := template.FuncMap{
		"initialLower": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			r := []rune(s)
			r[0] = []rune(strings.ToLower(string(r[0])))[0]
			return string(r)
		},
		"hasField": func(fieldName string) bool {
			// For testing, return true for common fields
			return fieldName == "Name" || fieldName == "Email" || fieldName == "Id" || 
				   fieldName == "CreatedAt" || fieldName == "UpdatedAt"
		},
	}

	// Parse template
	tmpl, err := template.New("usecase").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Create mock data for UserService with ORM information
	mockData := map[string]interface{}{
		"Service": struct {
			GoName  string
			Methods []struct {
				GoName string
				Input  struct {
					GoIdent struct {
						GoName string
					}
				}
				Output struct {
					GoIdent struct {
						GoName string
					}
				}
			}
		}{
			GoName: "UserService",
			Methods: []struct {
				GoName string
				Input  struct {
					GoIdent struct {
						GoName string
					}
				}
				Output struct {
					GoIdent struct {
						GoName string
					}
				}
			}{
				{
					GoName: "CreateUser",
					Input: struct {
						GoIdent struct {
							GoName string
						}
					}{
						GoIdent: struct {
							GoName string
						}{
							GoName: "CreateUserRequestDTO",
						},
					},
					Output: struct {
						GoIdent struct {
							GoName string
						}
					}{
						GoIdent: struct {
							GoName string
						}{
							GoName: "CreateUserResponseDTO",
						},
					},
				},
				{
					GoName: "GetUser",
					Input: struct {
						GoIdent struct {
							GoName string
						}
					}{
						GoIdent: struct {
							GoName string
						}{
							GoName: "GetUserRequestDTO",
						},
					},
					Output: struct {
						GoIdent struct {
							GoName string
						}
					}{
						GoIdent: struct {
							GoName string
						}{
							GoName: "GetUserResponse",
						},
					},
				},
			},
		},
		"ServiceNoSuf": "User",
		"Base":        "github.com/harryosmar/protobuf-go",
		"Filename":    "user_service_usecase.go",
		"Package":     "usecase",
		"ProtoPkg":    "github.com/harryosmar/protobuf-go/gen/user",
		"ServiceLower": "user_service",
		// ORM-related information
		"ORM": map[string]interface{}{
			"EntityName": "UserEntityORM",
			"DTOName":    "UserDTO",
			"IDType":     "uint32",
			"Fields": []string{
				"Id",
				"Name",
				"Email",
				"CreatedAt",
				"UpdatedAt",
			},
		},
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mockData)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	output := buf.String()

	// Verify expected content is present
	expectedStrings := []string{
		"package usecase",
		"type UserServiceUsecase interface",
		"CreateUser(context.Context, *pb.CreateUserRequestDTO) (*pb.CreateUserResponseDTO, error)",
		"GetUser(context.Context, *pb.GetUserRequestDTO) (*pb.GetUserResponse, error)",
		"type userServiceUsecase struct",
		"userServiceRepo repository.ServiceRepository[pb.UserEntityORM, uint32]",
		"func NewUserServiceUsecase(",
		"func (u *userServiceUsecase) ormToDTO(orm *pb.UserEntityORM) *pb.UserDTO",
		"Id: orm.Id,",
		"Name: orm.Name,",
		"Email: orm.Email,",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it didn't", expected)
		}
	}

	// Verify ORM-related content is present
	ormStrings := []string{
		"pb.UserEntityORM",
		"pb.UserDTO",
		"uint32",
		"repository.ServiceRepository[pb.UserEntityORM, uint32]",
	}

	for _, expected := range ormStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain ORM-related string '%s', but it didn't", expected)
		}
	}

	// Verify no reflection code is present
	unexpectedStrings := []string{
		"reflect.ValueOf",
		"reflect.New",
		"reflect.Type",
	}

	for _, unexpected := range unexpectedStrings {
		if strings.Contains(output, unexpected) {
			t.Errorf("Output should not contain '%s', but it does", unexpected)
		}
	}
}
