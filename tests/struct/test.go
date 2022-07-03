package _struct

//go:generate go run github.com/apitalist/enum/cmd/enum/ -type TestEnum
type TestEnum struct {
    value string
}

var (
    TestEnumA = TestEnum{"a"}
    TestEnumB = TestEnum{"b"}
    TestEnumC = TestEnum{"c"}
)
