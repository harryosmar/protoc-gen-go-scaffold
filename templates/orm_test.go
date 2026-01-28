package templates

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"
)

// Simple mock service for testing
type mockService struct {
	GoName  string
	Methods []mockMethod
}

type mockMethod struct {
	GoName string
	Input  mockGoIdent
	Output mockGoIdent
}

type mockGoIdent struct {
	GoIdent mockIdent
}

type mockIdent struct {
	GoName string
}

// Helper function to simulate initialLower template function
func initialLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]
	return string(r)
}

// Helper function to simulate hasField template function
func hasField(fieldName string) bool {
	// For testing, return true for common fields
	return fieldName == "Name" || fieldName == "Email" || fieldName == "Id" || 
	       fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "Description"
}

func TestUsecaseTemplateWithORM(t *testing.T) {
	// Load the template file
	tmplContent, err := os.ReadFile("usecase.tmpl")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	// Define template functions
	funcMap := template.FuncMap{
		"initialLower": initialLower,
		"hasField":     hasField,
	}

	// Parse template
	tmpl, err := template.New("usecase").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Create mock data for UserService with ORM information
	mockData := map[string]interface{}{
		"Service": mockService{
			GoName: "UserService",
			Methods: []mockMethod{
				{
					GoName: "CreateUser",
					Input:  mockGoIdent{GoIdent: mockIdent{GoName: "CreateUserRequestDTO"}},
					Output: mockGoIdent{GoIdent: mockIdent{GoName: "CreateUserResponseDTO"}},
				},
				{
					GoName: "GetUser",
					Input:  mockGoIdent{GoIdent: mockIdent{GoName: "GetUserRequestDTO"}},
					Output: mockGoIdent{GoIdent: mockIdent{GoName: "GetUserResponse"}},
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
