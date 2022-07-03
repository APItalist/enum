package string

import (
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
