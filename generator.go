package enum

import (
    "bytes"
    _ "embed"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "os"
    "sort"
    "strings"
    "text/template"
)

// New creates a new generator.
func New() Generator {
    return &generator{}
}

// Generator generates helper code for working with enums.
type Generator interface {
    // Generate generates the enum helpers, as well as test code if desired. If the generation fails, an error is
    // returned.
    Generate(spec Spec) ([]byte, []byte, error)
}

// Spec is the configuration for the Generator.
type Spec struct {
    // Type is the type to generate enum helpers for.
    Type string
    // Directory is the source directory for the generation.
    Directory string
    // GenerateTests generates test code for the generated enum code.
    GenerateTests bool
    // GenerateValidate generates a validation function.
    GenerateValidate bool
    // GenerateValues generates a values function.
    GenerateValues bool
    // GenerateListType generates a list type.
    GenerateListType bool
    // GenerateJSON generates a JSON unmarshaller with validation.
    GenerateJSON bool
}

// Validate validates the specification.
func (s *Spec) Validate() error {
    if s.Type == "" {
        return fmt.Errorf("type cannot be empty")
    }
    if s.Directory == "" {
        var err error
        s.Directory, err = os.Getwd()
        if err != nil {
            return fmt.Errorf("no source directory set and getting current working directory failed (%w)", err)
        }
    }
    return nil
}

// TemplateScope is the data for rendering templates.
type TemplateScope struct {
    // Exported is true if the type is exported.
    Exported bool
    // Package is the package name for the type.
    Package string
    // First is the first letter of the type.
    First string
    // LowerFirst is the first letter of the type, lower case.
    LowerFirst string
    // Type is the name of the type being generated.
    Type string
    // Values is a list of values for the enum.
    Values []string
    // ConvertToString is the prefix for converting to string.
    ConvertToString string
    // ConverToStringEnd is the suffix for converting to string.
    ConvertToStringEnd string
    // RawType is the underlying type.
    RawType string
    // Spec is the input spec.
    Spec Spec
    // ConvertImports holds a list of imports used for conversion.
    ConvertImports map[string]struct{}
}

func (t TemplateScope) Imports(added ...string) []string {
    imports := make(map[string]struct{}, len(t.ConvertImports))
    for k, v := range t.ConvertImports {
        imports[k] = v
    }
    for _, a := range added {
        imports[a] = struct{}{}
    }
    result := make([]string, len(imports))
    i := 0
    for imp := range imports {
        result[i] = imp
        i++
    }
    sort.SliceStable(result, func(i, j int) bool {
        return strings.Compare(result[i], result[j]) < 0
    })
    return result
}

//go:embed enum.go.tpl
var enumTemplate string

//go:embed enum_test.go.tpl
var enumTestTemplate string

type generator struct {
}

func (g generator) Generate(spec Spec) ([]byte, []byte, error) {
    if err := spec.Validate(); err != nil {
        return nil, nil, err
    }
    tplData, err := g.parse(spec)
    if err != nil {
        return nil, nil, fmt.Errorf("parse failed (%w)", err)
    }
    switch tplData.RawType {
    case "string":
        tplData.ConvertToString = `string(`
        tplData.ConvertToStringEnd = `)`
    default:
        tplData.ConvertToString = `fmt.Sprintf("%v", `
        tplData.ConvertToStringEnd = ")"
        tplData.ConvertImports = map[string]struct{}{"fmt": {}}
    }
    enumResult, err := g.renderTemplate(enumTemplate, tplData)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to render enum.go.tpl template (%w)", err)
    }
    var enumTestResult []byte
    if spec.GenerateTests {
        enumTestResult, err = g.renderTemplate(enumTestTemplate, tplData)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to render enum.go.tpl template (%w)", err)
        }
    }
    return enumResult, enumTestResult, nil
}

func (g generator) parse(s Spec) (*TemplateScope, error) {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, s.Directory, nil, parser.SkipObjectResolution)
    if err != nil {
        return nil, fmt.Errorf("failed to parse %s (%w)", s.Directory, err)
    }

    var typePkg *ast.Package
    var values []*ast.ValueSpec
    var typeSpec *ast.TypeSpec
    for _, pkg := range pkgs {
        for _, f := range pkg.Files {
            for _, decl := range f.Decls {
                if genDecl, ok := decl.(*ast.GenDecl); ok {
                    switch genDecl.Tok {
                    case token.TYPE:
                        typeDecl := genDecl.Specs[0].(*ast.TypeSpec)
                        if typeDecl.Name.Name == s.Type {
                            typeSpec = typeDecl
                        }
                        typePkg = pkg
                    case token.CONST:
                        iota := false
                        for _, spec := range genDecl.Specs {
                            valueSpec, ok := spec.(*ast.ValueSpec)
                            if !ok {
                                continue
                            }
                            if len(valueSpec.Values) == 1 {
                                value := valueSpec.Values[0]
                                if ident, ok := value.(*ast.Ident); ok && ident.Name == "iota" {
                                    iota = true
                                } else {
                                    iota = false
                                }
                            }
                            ident, ok := valueSpec.Type.(*ast.Ident)
                            if !ok && !iota {
                                continue
                            }
                            if !iota && (ident == nil || ident.Name != s.Type) {
                                continue
                            }
                            values = append(values, valueSpec)
                        }
                    case token.VAR:
                        for _, spec := range genDecl.Specs {
                            valueSpec := spec.(*ast.ValueSpec)
                            if len(valueSpec.Values) != 1 {
                                continue
                            }
                            compositeLit, ok := valueSpec.Values[0].(*ast.CompositeLit)
                            if !ok {
                                continue
                            }
                            ident, ok := compositeLit.Type.(*ast.Ident)
                            if !ok {
                                continue
                            }
                            if ident.Name != s.Type {
                                continue
                            }
                            values = append(values, valueSpec)
                        }
                    }
                }
            }
        }
    }
    if typeSpec == nil {
        return nil, fmt.Errorf("type %s not found in %s", s.Type, s.Directory)
    }
    if len(values) == 0 {
        return nil, fmt.Errorf("no const declarations found for type %s", s.Type)
    }
    tplValues := []string{}
    for _, v := range values {
        for _, name := range v.Names {
            tplValues = append(tplValues, name.Name)
        }
    }
    rawType := ""
    if ident, ok := typeSpec.Type.(*ast.Ident); ok {
        rawType = ident.Name
    }
    tplData := &TemplateScope{
        Exported:       typeSpec.Name.IsExported(),
        Package:        typePkg.Name,
        First:          strings.ToLower(s.Type[0:1]),
        LowerFirst:     strings.ToLower(s.Type[0:1]) + s.Type[1:],
        Type:           s.Type,
        Values:         tplValues,
        RawType:        rawType,
        Spec:           s,
        ConvertImports: map[string]struct{}{},
    }
    return tplData, nil
}

func (g generator) renderTemplate(templateSource string, tplData *TemplateScope) ([]byte, error) {
    tpl := template.Must(template.New("enum").Parse(templateSource))
    data := &bytes.Buffer{}
    if err := tpl.Execute(data, tplData); err != nil {
        log.Fatal(err)
    }
    return data.Bytes(), nil
}
