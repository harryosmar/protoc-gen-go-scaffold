package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/harryosmar/protoc-gen-go-scaffold/templates/helpers"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		NewScaffoldModule(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

// ScaffoldModule implements the protoc-gen-star Module interface to generate code scaffolding
type ScaffoldModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

// NewScaffoldModule returns a new ScaffoldModule
func NewScaffoldModule() *ScaffoldModule {
	return &ScaffoldModule{
		ModuleBase: &pgs.ModuleBase{},
	}
}

// Name returns the name of this module
func (m *ScaffoldModule) Name() string {
	return "scaffold"
}

// GetPkg returns the package name for the associated file
func (m *ScaffoldModule) GetPkg(f pgs.File) string {
	return m.ctx.ImportPath(f).String()
}

// InitContext initializes the module with the Go context
func (m *ScaffoldModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
}

// Execute generates the code for each service defined in the proto files
func (m *ScaffoldModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		if len(f.Services()) == 0 {
			continue
		}
		m.generateServerFiles(f)
		m.generateUsecaseFiles(f)
		m.generateRepositoryFiles(f)
		m.generateErrorFiles(f)
		m.generateMainFile(f)
	}
	return m.ModuleBase.Artifacts()
}

// processTemplate processes a template file and returns the rendered content
func (m *ScaffoldModule) processTemplate(name, tmplPath string, data map[string]interface{}) (string, error) {
	// Check if a templates_dir parameter was provided
	templatesDir := m.Parameters().Str("templates_dir")
	if templatesDir == "" {
		templatesDir = "."
	}

	// Get base package parameter or use default
	basePackage := m.Parameters().Str("base")
	if basePackage == "" {
		basePackage = "github.com/harryosmar/protobuf-go"
	}

	// Add Base to the data map
	data["Base"] = basePackage

	// Use path from templates_dir if provided
	fullPath := filepath.Join(templatesDir, tmplPath)

	tmplData, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	// Parse the template with helper functions from templates/helpers
	tmpl := template.New(name)
	tmpl = helpers.AddTemplateFuncs(tmpl)
	tmpl, err = tmpl.Parse(string(tmplData))

	if err != nil {
		return "", err
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// serviceNoSuffix removes "Service" suffix and returns properly cased version
func serviceNoSuffix(s string) string {
	// Remove "Service" suffix
	name := strings.TrimSuffix(s, "Service")
	// Ensure first letter is uppercase for consistent casing
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return name
}

// getServiceBaseName removes "Service" suffix from a service name to avoid duplication
func getServiceBaseName(serviceName string) string {
	return serviceNoSuffix(serviceName)
}

// getServiceBaseLower returns a lowercase version of the service name without "Service" suffix
func getServiceBaseLower(serviceName string) string {
	return strings.ToLower(getServiceBaseName(serviceName))
}

// generateServerFiles generates the server layer code
func (m *ScaffoldModule) generateServerFiles(f pgs.File) {
	for _, service := range f.Services() {
		serviceName := service.Name().String()
		serviceLowerBase := strings.ToLower(serviceNoSuffix(serviceName))

		// Prepare template data - pass Service object and processed variables for backward compatibility
		templateData := map[string]interface{}{
			"Service":      service,
			"ServiceName":  serviceNoSuffix(serviceName),
			"ServiceLower": serviceLowerBase,
			"File":         f,
			"Import":       m.ctx.ImportPath(f),
		}

		// Process template with Service object and processed variables for backward compatibility

		// Process template - file name generation handled in template
		fileName := strings.ToLower(serviceNoSuffix(serviceName)) + "_server.go"
		content, err := m.processTemplate(fileName, "templates/server.tmpl", templateData)
		if err != nil {
			m.ModuleBase.Logf("Error processing template: %v", err)
			continue
		}

		// Add the generated file to the server directory
		m.ModuleBase.AddGeneratorFile("server/"+fileName, content)
	}
}

// generateUsecaseFiles generates the usecase layer code
func (m *ScaffoldModule) generateUsecaseFiles(f pgs.File) {
	for _, service := range f.Services() {
		serviceName := service.Name().String()
		serviceLowerBase := strings.ToLower(serviceNoSuffix(serviceName))

		// Prepare template data - pass Service object and processed variables for backward compatibility
		templateData := map[string]interface{}{
			"Service":      service,
			"ServiceName":  serviceNoSuffix(serviceName),
			"ServiceLower": serviceLowerBase,
			"File":         f,
			"Import":       m.ctx.ImportPath(f),
		}

		// Process template - file name generation handled in template
		fileName := strings.ToLower(serviceNoSuffix(serviceName)) + "_usecase.go"
		content, err := m.processTemplate(fileName, "templates/usecase.tmpl", templateData)
		if err != nil {
			m.ModuleBase.Logf("Error processing template: %v", err)
			continue
		}

		// Add the generated file to the usecase directory
		m.ModuleBase.AddGeneratorFile("usecase/"+fileName, content)
	}
}

// generateRepositoryFiles generates the repository layer code
func (m *ScaffoldModule) generateRepositoryFiles(f pgs.File) {
	for _, service := range f.Services() {
		serviceName := service.Name().String()
		serviceLowerBase := strings.ToLower(serviceNoSuffix(serviceName))

		// Prepare template data - pass Service object and processed variables for backward compatibility
		templateData := map[string]interface{}{
			"Service":      service,
			"ServiceName":  serviceNoSuffix(serviceName),
			"ServiceLower": serviceLowerBase,
			"File":         f,
			"Import":       m.ctx.ImportPath(f),
		}

		// Process template - file name generation handled in template
		fileName := strings.ToLower(serviceNoSuffix(serviceName)) + "_repository.go"
		content, err := m.processTemplate(fileName, "templates/repository.tmpl", templateData)
		if err != nil {
			m.ModuleBase.Logf("Error processing template: %v", err)
			continue
		}

		// Add the generated file to the repository directory
		m.ModuleBase.AddGeneratorFile("repository/"+fileName, content)
	}
}

// generateErrorFiles generates the error code files
func (m *ScaffoldModule) generateErrorFiles(f pgs.File) {
	for _, service := range f.Services() {
		serviceName := service.Name().String()
		serviceLowerBase := strings.ToLower(serviceNoSuffix(serviceName))

		// Prepare template data - pass Service object and processed variables for backward compatibility
		templateData := map[string]interface{}{
			"Service":      service,
			"ServiceName":  serviceNoSuffix(serviceName),
			"ServiceLower": serviceLowerBase,
			"File":         f,
			"Import":       m.ctx.ImportPath(f),
		}

		// Process template - file name generation handled in template
		// Reuse the existing serviceName variable
		fileName := strings.ToLower(serviceNoSuffix(serviceName)) + "_codes.go"
		content, err := m.processTemplate(fileName, "templates/error_code.tmpl", templateData)
		if err != nil {
			m.ModuleBase.Logf("Error processing template: %v", err)
			continue
		}

		// Add the generated file to the error directory
		m.ModuleBase.AddGeneratorFile("error/"+fileName, content)
	}
}

// generateMainFile generates the main.go file with init function
func (m *ScaffoldModule) generateMainFile(f pgs.File) {
	for _, service := range f.Services() {
		serviceName := service.Name().String()
		serviceLowerBase := strings.ToLower(serviceNoSuffix(serviceName))

		// Prepare template data - pass Service object and processed variables for backward compatibility
		templateData := map[string]interface{}{
			"Service":      service,
			"ServiceName":  serviceNoSuffix(serviceName),
			"ServiceLower": serviceLowerBase,
			"File":         f,
			"Import":       m.ctx.ImportPath(f),
		}

		// Process template
		content, err := m.processTemplate("main.go", "templates/main.tmpl", templateData)
		if err != nil {
			m.ModuleBase.Logf("Error processing template: %v", err)
			continue
		}

		// Add the generated main file
		m.ModuleBase.AddGeneratorFile("main.go", content)
	}
}

// Parameters returns the module's current parameters
func (m *ScaffoldModule) Parameters() pgs.Parameters {
	return m.BuildContext.Parameters()
}
