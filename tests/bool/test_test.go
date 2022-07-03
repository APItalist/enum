package string

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    if _, err := os.Stat("enum_TestEnum.go"); err != nil {
        log.Fatalf("enum_TestEnum.go not found, please run go generate")
    }
    m.Run()
}

func TestValue(t *testing.T) {
    for _, v := range []TestEnum{
        TestEnumA, TestEnumB,
    } {
        t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
            if err := v.Validate(); err != nil {
                t.Fatalf("%v failed validation (%v)", v, err)
            }
        })
    }
}

func TestJSON(t *testing.T) {
    for _, v := range []TestEnum{
        TestEnumA, TestEnumB,
    } {
        t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
            data, err := json.Marshal(v)
            if err != nil {
                t.Fatalf("failed to marshal %v (%v)", v, err)
            }
            var result TestEnum
            if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("failed to unmashal %v (%v)", v, err)
            }
            if result != v {
                t.Fatalf("mismatch (expected: %v, got: %v)", v, result)
            }
        })
    }
}
