package string

//go:generate go run github.com/apitalist/enum/cmd/enum/ -type TestEnum
type TestEnum string

const (
    TestEnumA TestEnum = "a"
    TestEnumB TestEnum = "b"
    TestEnumC TestEnum = "c"
)
