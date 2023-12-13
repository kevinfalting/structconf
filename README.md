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

### Secrets

The builtin secret handler expects base64 encoded RSA strings, prefixed with `secret://`, to look like `secret://<some_base64_encoded_rsa_encrypted_ciphertext>`.

The advantage of this is that handlers can read in encrypted or decrypted values, and once they get to the secret handler, it will either skip it or decrypt a value prefixed with `secret://`, since you may want to provide a decrypted value as an environment variable or flag. You can also provide an encrypted value with that prefix as an environment variable or flag, and it will be decrypted using the provided key.

## Supporting Unsupported Types

The parser will prioritize value fields that satisfy the `encoding.TextUnmarshaler` or `encoding.BinaryUnmarshaler`, in that order. If you need to support an unsupported type like a map or slice, then create a user defined type that satisfies either interface.

## Module Opinions
- Configuration should only use value semantics. There are ways around this in this module, but pointers and reference types should generally be avoided. Configuration tends to be shared across goroutines.
- A field should not need to be aware of it's position in nested layers of structs. The field should contain all the information required to lookup its value.
