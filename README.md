# structconf

An extensible module for struct configuration with sensible defaults.

## Module Sate

`structconf` is currently in v0. There are possible breaking changes.

##  Usage

Built in handlers support the following comma separated struct tags defined inside of the main struct tag key `conf`.

| tag | usage | description |
|-|-|-|
| `env` | `env:APP_NAME` | defines the environment variable the environment variable handler uses to lookup the value. |
| `flag` | `flag:app-name` | defines the command line flag to lookup the value. The flag handler is optional. |
| `usage` | `usage:this is how you use it` | defines the usage text in the help message when using the flag handler. |
| `default` | `default:the app name` | defines the default value for the field. |
| `required` | `required` | defines whether the field is required or not. No value necessary. |

```go
type Config struct {
    Name `conf:"env:APP_NAME,flag:app-name,usage:this is how you use it,default:the app name,required"`
}
```

> Using the `default` and `required` tags together won't cause any errors, although they may be redundant.

The precedence of the default configuration is applied in the following order:

1. Default Value (defined by the `default` tag)
1. Field Value (when an initialized, non-zero value is present in the provided struct)
1. Environment Variable (defined by the `env` tag)
1. Command Line Flag (defined by the `flag` tag, when the flag handler is enabled)

## Supporting Unsupported Types

The parser will prioritize value fields that satisfy the `encoding.TextUnmarshaler` or `encoding.BinaryUnmarshaler`, in that order. If you need to support an unsupported type like a map or slice, then create a user defined type that satisfies either interface.

## Module Structure

`structconf` is separated into three distinct packages. The primary reason is to manage the api for each of the packages, keeping them tightly focused.

#### `github.com/kevinfalting/structconf`

This is the module name and top level package, the api here is the first thing encountered and should be kept very simple, focused only on what this module is about: configuring structs, by using "sensible" defaults. The api here is only what is needed to quickly parse a struct without needing to do any extra setup.

This package is only an abstraction layer of the other two packages in this module.

#### `github.com/kevinfalting/structconf/stronf`

This package makes up the core of this module, and contains all of the opinions about what fields are considered settable, what interfaces to respect, and how to coerce types into a type and value that is settable. It defines what the function signature is for a handler, and models a struct field that can be used in handlers to determine what value to set.

These are publicly exposed to encourage building wrappers around this package for any _boutique_ configuration environments.

#### `github.com/kevinfalting/structconf/confhandler`

This package contains a set of handlers that can be used out-of-the-box. There's a small amount of handlers here, but they should be enough for the typical application, or show how to build custom ones. A handler is just a function, so it's simple enough to make your own.

## Extending `structconf`

There's a limited set of handlers in this module, partly because more haven't been built yet, and partly because this should be kept to only the standard library. There are many great implementations of various file parsers and remote configuration management, but I didn't want to import them here. `structconf` exposes everything needed to write custom handlers that can perform what is needed in a specialized environment.

## Module Philosophies
- Configuration should only use value semantics. There are ways around this in this module, but pointers and reference types should generally be avoided. Configuration tends to be shared across goroutines.
- A field should not need to be aware of it's position in nested layers of structs. The field should contain all the information required to lookup its value.
- This module makes no assumptions about what names to use for environment variables, flags, or anything else. Everything must be explicitly provided, otherwise it is ignored.
- Fields that depend on other fields for their value are not supported in this module. That is better suited as a method which builds the value after it's been parsed by this module.
- No dependencies outside of the standard library.
