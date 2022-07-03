# Better enums for Go

This library provides automatic enum code generation for Go.

## Installation

Install by typing:

```
go install github.com/apitalist/enum/cmd/enum
```

## Usage

First, add your enum:

```go
type MyEnum string

const (
    MyEnumA MyEnum = "a"
    MyEnumB MyEnum = "b"
)
```

Then add a generate line:

```go
//go:generate go run github.com/apitalist/enum/cmd/enum/ -type MyEnum
```

This will generate a file called `enum_MyEnum.go` as well as `enum_MyEnum_test.go` with the following functions:

- `Validate()` validates the enum and returns an error if the value is not one of the specified constants.
- `MyEnumValues()` returns a list of valid value for the enum.
- `MyEnums` is a type for a list of `MyEnum` values.
- `MyEnumValueStrings()` returns a list of valid values as strings for the enum.

It also supports other enum types:

```go
//go:generate go run github.com/apitalist/enum/cmd/enum/ -type MyEnum
type MyEnum int

const (
    MyEnumA MyEnum = iota
    MyEnumB
    MyEnumC
)
```

Or safe enums. You may want to implement JSON and text marshalling though.

```go
//go:generate go run github.com/apitalist/enum/cmd/enum/ -type MyEnum
type MyEnum struct {
    value string
}

var (
    MyEnumA MyEnum = MyEnum{"a"}
    MyEnumB MyEnum = MyEnum{"b"}
    MyEnumC MyEnum = MyEnum{"c"}
)
```

## Options

You can pass the following options to the enum generator:

| Option        | Default value      | Description                                                                |
|---------------|--------------------|----------------------------------------------------------------------------|
| `-type`       | *none*             | Type to generate enum from                                                 |
| `-source`     | `.`                | Source directory                                                           |
| `-target`     | `enum_TYPENAME.go` | Target file name. The tests will be generated into a file named `_test.go` |
| `-notests`    | `false`            | Do not generate tests                                                      |
| `-novalidate` | `false`            | Do not generate `Validate()` functions                                     |
| `-novalues`   | `false`            | Do not generate `Values()` functions                                       |
| `-nolist`     | `false`            | Do not generate list types                                                 |

## Documentation

For API documentation please see [pkg.go.dev/github.com/apitalist/enum](https://pkg.go.dev/github.com/apitalist/enum).

## License

APItalist is licensed under [the Apache 2.0 license](LICENSE).