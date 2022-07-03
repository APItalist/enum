{{- /*gotype: github.com/apitalist/enum.TemplateScope */ -}}
// Code generated by generate-enum. DO NOT EDIT.

package {{ .Package }}

import (
{{ range .Imports "testing"}}    "{{ . }}"
{{ end}})

{{ if .Spec.GenerateValidate }}
// Test{{ .Type }}Validate tests the Validate method of {{ .Type }}.
func Test{{ .Type }}Validate(t *testing.T) {
    for _, v := range {{ .Type }}Values() {
        t.Run({{.ConvertToString}}v{{ .ConvertToStringEnd }}, func(t *testing.T) {
            if err := v.Validate(); err != nil {
                t.Fatalf("%v failed validation (%v)", v, err)
            }
        })
    }
}

// Test{{ .Type }}sValidate tests the Validate method of {{ .Type }}s.
func Test{{ .Type }}sValidate(t *testing.T) {
    if err := {{ .Type }}Values().Validate(); err != nil {
        t.Fatalf("failed validation (%v)", err)
    }
}
{{- end }}