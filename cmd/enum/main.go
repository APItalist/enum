package main

import (
    "flag"
    "io/ioutil"
    "log"
    "strings"

    "github.com/apitalist/enum"
)

func main() {
    spec := enum.Spec{
        Type:             "",
        Directory:        "",
        GenerateTests:    false,
        GenerateValidate: false,
        GenerateValues:   false,
        GenerateListType: false,
    }
    var noGenerateTests bool
    var noGenerateValidate bool
    var noGenerateValues bool
    var noGenerateListType bool
    targetFile := ""
    flag.StringVar(&spec.Type, "type", spec.Type, "Type to generate enum helpers for.")
    flag.StringVar(&spec.Directory, "source", spec.Directory, "Source directory for the type.")
    flag.StringVar(&targetFile, "target", targetFile, "Target file to write to. Tests will be written to a separate test file.")
    flag.BoolVar(&noGenerateTests, "notests", noGenerateTests, "Do not generate test code.")
    flag.BoolVar(&noGenerateValidate, "novalidate", noGenerateValidate, "Do not generate validation functions.")
    flag.BoolVar(&noGenerateValues, "novalues", noGenerateValues, "Do not generate values functions.")
    flag.BoolVar(&noGenerateListType, "nolist", noGenerateListType, "Do not generate list type.")
    flag.Parse()

    if targetFile == "" {
        targetFile = "enum_" + spec.Type + ".go"
    }
    spec.GenerateValidate = !noGenerateValidate
    spec.GenerateValues = !noGenerateValues
    spec.GenerateListType = !noGenerateListType
    spec.GenerateTests = !noGenerateTests

    generator := enum.New()
    enumContent, testContent, err := generator.Generate(spec)
    if err != nil {
        log.Fatal(err)
    }

    if err := ioutil.WriteFile(targetFile, enumContent, 0644); err != nil {
        log.Fatalf("failed to write target file %s (%v)", targetFile, err)
    }
    if noGenerateTests || len(testContent) == 0 {
        return
    }

    testTargetFile := strings.Replace(targetFile, ".go", "_test.go", 1)
    if err := ioutil.WriteFile(
        testTargetFile,
        testContent,
        0644,
    ); err != nil {
        log.Fatalf("failed to write target test file %s (%v)", testTargetFile, err)
    }
}
