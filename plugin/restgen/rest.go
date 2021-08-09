package restgen

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/99designs/gqlgen/internal/code"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

func New(filename string, typename string) plugin.Plugin {
	return &Plugin{filename: filename, typeName: typename}
}

type Plugin struct {
	filename string
	typeName string
}

var _ plugin.CodeGenerator = &Plugin{}
var _ plugin.ConfigMutator = &Plugin{}

func (m *Plugin) Name() string {
	return "restgen"
}

func (m *Plugin) MutateConfig(cfg *config.Config) error {
	_ = syscall.Unlink(m.filename)
	return nil
}

// ADE:
func DbgPrint(data *codegen.Data, object *codegen.Object) {
	// data.Objects
	// data.QueryRoot.Fields
	// _ = data.Objects[0].Fields[0].ShortResolverDeclaration()
	// _ = data.Objects[0].Fields[0].Arguments[0].Name
	fmt.Printf("\n=> objects: %#v\n", object.Type)
	if object == nil {
		return
	}

	for _, field := range object.Fields {
		fmt.Printf("=> fields: %#v, %s\n", field, field.Name)
		if field == nil {
			return
		}

		fmt.Println("=> field: ", field.Name, field.FieldDefinition.Name)

		for _, innerField := range field.TypeReference.Definition.Fields {
			fmt.Println(innerField.Name, innerField.Type)
			innerObject := data.Objects.ByName(innerField.Name)
			if innerObject != nil {
				DbgPrint(data, innerObject)
			}
		}
	}
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	DbgPrint(data, data.Objects.ByName("query"))

	abs, err := filepath.Abs(m.filename)
	if err != nil {
		return err
	}
	pkgName := code.NameForDir(filepath.Dir(abs))

	return templates.Render(templates.Options{
		PackageName: pkgName,
		Filename:    m.filename,
		Data: &ResolverBuild{
			Data:     data,
			TypeName: m.typeName,
		},
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

type ResolverBuild struct {
	*codegen.Data

	TypeName string
}
