# structconf

An extensible module for struct configuration with sensible defaults.

## Module Sate

`structconf` is currently in progress and not ready for a v1 tag. It's being tested in various scenarios and refined as we learn more and bugs are uncovered.

Expect breaking changes.

##  Usage

Built in handlers support the following comma separated struct tags defined inside of the main struct tag key `conf`.

| tag | usage | description |
|-|-|-|
| `env` | `env:APP_NAME` | defines the environment variable the environment variable handler uses to lookup the value. |
| `flag` | `flag:app-name` | defines the command line flag to lookup the value. |
| `default` | `default:the app name` | defines the default value for the field. |
| `required` | `required` | defines whether the field is required or not. No value necessary. |

```go
type Config struct {
    Name `conf:"env:APP_NAME,flag:app-name,default:the app name,required"`
}
```

> Using the `default` and `required` tags together won't cause any errors, although they may be redundant.

The precedence of the default configuration is applied in the following order:

1. Default Value
1. Field Value
1. Environment Variable
1. Command Line Flag

## Supporting Unsupported Types

The parser will prioritize value fields that satisfy the `encoding.TextUnmarshaler` or `encoding.BinaryUnmarshaler`, in that order. If you need to support an unsupported type like a map or slice, then create a user defined type that satisfies either interface.

## Module Opinions
- Configuration should only use value semantics. There are ways around this in this module, but pointers and reference types should generally be avoided. Configuration tends to be shared across goroutines.
- A field should not need to be aware of it's position in nested layers of structs. The field should contain all the information required to lookup its value.
