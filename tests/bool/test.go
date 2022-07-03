package string

//go:generate go run github.com/apitalist/enum/cmd/enum/ -type TestEnum
type TestEnum bool

const (
    TestEnumA TestEnum = true
    TestEnumB TestEnum = false
)
