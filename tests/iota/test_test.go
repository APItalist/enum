package string

import (
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
        TestEnumA, TestEnumB, TestEnumC,
    } {
        t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
            if err := v.Validate(); err != nil {
                t.Fatalf("%v failed validation (%v)", v, err)
            }
        })
    }
}
