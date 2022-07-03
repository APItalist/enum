package string

//go:generate go run github.com/apitalist/enum/cmd/enum/ -type TestEnum
type TestEnum int

const (
    TestEnumA TestEnum = iota
    TestEnumB TestEnum = iota
    TestEnumC TestEnum = iota
)
