package generator

import (
	"fmt"
	"os"
	"text/template"
)

// CodeGenerator generates code based on entity definitions
type CodeGenerator struct {
	enabled bool
}

// NewCodeGenerator creates a new CodeGenerator instance
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		enabled: true,
	}
}

// GenerateHelper generates helper code for an entity
func (cg *CodeGenerator) GenerateHelper(entityName string) error {
	if !cg.enabled {
		return nil
	}

	// Helper template
	helperTemplate := `// Package helper implements the {{.EntityName}} helper for the OIDC application.
package helper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// {{.EntityName}}Helper defines the {{.EntityNameLower}} helper interface
type {{.EntityName}}Helper interface {
	IBaseHelper

	// Validate{{.EntityName}} validates {{.EntityNameLower}} entity
	Validate{{.EntityName}}({{.EntityNameLower}} *model.{{.EntityName}}) error

	// Format{{.EntityName}} formats {{.EntityNameLower}} entity
	Format{{.EntityName}}({{.EntityNameLower}} *model.{{.EntityName}}) *model.{{.EntityName}}
}

// {{.EntityNameLower}}Helper implements {{.EntityName}}Helper interface
type {{.EntityNameLower}}Helper struct {
	version string
}

// New{{.EntityName}}Helper creates a new {{.EntityName}}Helper instance
func New{{.EntityName}}Helper() {{.EntityName}}Helper {
	return &{{.EntityNameLower}}Helper{
		version: "1.0.0",
	}
}

// Validate validates the entity
func (uh *{{.EntityNameLower}}Helper) Validate(entity interface{}) error {
	if {{.EntityNameLower}}, ok := entity.(*model.{{.EntityName}}); ok {
		return uh.Validate{{.EntityName}}({{.EntityNameLower}})
	}
	return fmt.Errorf("invalid entity type")
}

// Format formats the entity
func (uh *{{.EntityNameLower}}Helper) Format(entity interface{}) interface{} {
	if {{.EntityNameLower}}, ok := entity.(*model.{{.EntityName}}); ok {
		return uh.Format{{.EntityName}}({{.EntityNameLower}})
	}
	return entity
}

// Validate{{.EntityName}} validates {{.EntityNameLower}} entity
func (uh *{{.EntityNameLower}}Helper) Validate{{.EntityName}}({{.EntityNameLower}} *model.{{.EntityName}}) error {
	// Implement validation logic
	return nil
}

// Format{{.EntityName}} formats {{.EntityNameLower}} entity
func (uh *{{.EntityNameLower}}Helper) Format{{.EntityName}}({{.EntityNameLower}} *model.{{.EntityName}}) *model.{{.EntityName}} {
	// Implement formatting logic
	return {{.EntityNameLower}}
}

// HealthCheck checks the health of the helper
func (uh *{{.EntityNameLower}}Helper) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the helper
func (uh *{{.EntityNameLower}}Helper) GetVersion() string {
	return uh.version
}
`

	// Parse template
	tmpl, err := template.New("helper").Parse(helperTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse helper template: %w", err)
	}

	// Create output file
	fileName := fmt.Sprintf("internal/helper/%s_helper_impl.go", entityName)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create helper file: %w", err)
	}
	defer file.Close()

	// Execute template
	data := map[string]interface{}{
		"EntityName":      entityName,
		"EntityNameLower": toLowerFirst(entityName),
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute helper template: %w", err)
	}

	return nil
}

// GenerateMapper generates mapper code for an entity
func (cg *CodeGenerator) GenerateMapper(entityName string) error {
	if !cg.enabled {
		return nil
	}

	// Mapper template
	mapperTemplate := `// Package mapper implements the {{.EntityName}} mapper for the OIDC application.
package mapper

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Full-finger/OIDC/internal/model"
	_ "github.com/lib/pq"
)

// {{.EntityName}}Mapper defines the {{.EntityNameLower}} mapper interface
type {{.EntityName}}Mapper interface {
	IBaseMapper

	// Create{{.EntityName}} creates a new {{.EntityNameLower}}
	Create{{.EntityName}}(ctx context.Context, {{.EntityNameLower}} *model.{{.EntityName}}) error

	// Get{{.EntityName}}ByID gets a {{.EntityNameLower}} by ID
	Get{{.EntityName}}ByID(ctx context.Context, id int64) (*model.{{.EntityName}}, error)

	// Update{{.EntityName}} updates a {{.EntityNameLower}}
	Update{{.EntityName}}(ctx context.Context, {{.EntityNameLower}} *model.{{.EntityName}}) error

	// Delete{{.EntityName}} deletes a {{.EntityNameLower}}
	Delete{{.EntityName}}(ctx context.Context, id int64) error

	// List{{.EntityNamePlural}} lists {{.EntityNameLower}}s with pagination
	List{{.EntityNamePlural}}(ctx context.Context, offset, limit int) ([]*model.{{.EntityName}}, error)

	// Count{{.EntityNamePlural}} counts the total number of {{.EntityNameLower}}s
	Count{{.EntityNamePlural}}(ctx context.Context) (int64, error)
}

// {{.EntityNameLower}}Mapper implements {{.EntityName}}Mapper interface
type {{.EntityNameLower}}Mapper struct {
	db      *sql.DB
	version string
}

// New{{.EntityName}}Mapper creates a new {{.EntityName}}Mapper instance
func New{{.EntityName}}Mapper(db *sql.DB) {{.EntityName}}Mapper {
	return &{{.EntityNameLower}}Mapper{
		db:      db,
		version: "1.0.0",
	}
}

// Create{{.EntityName}} creates a new {{.EntityNameLower}}
func (um *{{.EntityNameLower}}Mapper) Create{{.EntityName}}(ctx context.Context, {{.EntityNameLower}} *model.{{.EntityName}}) error {
	// Implement create logic
	return nil
}

// Get{{.EntityName}}ByID gets a {{.EntityNameLower}} by ID
func (um *{{.EntityNameLower}}Mapper) Get{{.EntityName}}ByID(ctx context.Context, id int64) (*model.{{.EntityName}}, error) {
	// Implement get by ID logic
	return nil, nil
}

// Update{{.EntityName}} updates a {{.EntityNameLower}}
func (um *{{.EntityNameLower}}Mapper) Update{{.EntityName}}(ctx context.Context, {{.EntityNameLower}} *model.{{.EntityName}}) error {
	// Implement update logic
	return nil
}

// Delete{{.EntityName}} deletes a {{.EntityNameLower}}
func (um *{{.EntityNameLower}}Mapper) Delete{{.EntityName}}(ctx context.Context, id int64) error {
	// Implement delete logic
	return nil
}

// List{{.EntityNamePlural}} lists {{.EntityNameLower}}s with pagination
func (um *{{.EntityNameLower}}Mapper) List{{.EntityNamePlural}}(ctx context.Context, offset, limit int) ([]*model.{{.EntityName}}, error) {
	// Implement list logic
	return nil, nil
}

// Count{{.EntityNamePlural}} counts the total number of {{.EntityNameLower}}s
func (um *{{.EntityNameLower}}Mapper) Count{{.EntityNamePlural}}(ctx context.Context) (int64, error) {
	// Implement count logic
	return 0, nil
}

// HealthCheck checks the health of the mapper
func (um *{{.EntityNameLower}}Mapper) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the mapper
func (um *{{.EntityNameLower}}Mapper) GetVersion() string {
	return um.version
}
`

	// Parse template
	tmpl, err := template.New("mapper").Parse(mapperTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse mapper template: %w", err)
	}

	// Create output file
	fileName := fmt.Sprintf("internal/mapper/%s_mapper_impl.go", entityName)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create mapper file: %w", err)
	}
	defer file.Close()

	// Execute template
	data := map[string]interface{}{
		"EntityName":      entityName,
		"EntityNameLower": toLowerFirst(entityName),
		"EntityNamePlural": toPlural(entityName),
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute mapper template: %w", err)
	}

	return nil
}

// Enable enables the generator
func (cg *CodeGenerator) Enable() {
	cg.enabled = true
}

// Disable disables the generator
func (cg *CodeGenerator) Disable() {
	cg.enabled = false
}

// Helper functions for templates
func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]+32) + s[1:]
}

func toPlural(s string) string {
	// Simple pluralization - just add 's'
	return s + "s"
}
