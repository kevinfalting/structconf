/*
Package stronf exposes the core of the structconf module to support writing
custom configuration handlers.

The core ideas revolve around the [Field], [HandleFunc], and how values are
[Coerce]'d into the [Field].

The core properties of a [Field] are:
  - Not mutable from outside of the stronf package.
  - Expose only enough information to make programmatic decisions about what
    values to return.
  - Expose [reflect] package members only when necessary, and to avoid an
    unecessary abstraction layer.

The core properties of a settable field in [SettableFields] are:
  - Is settable from the [reflect] package's perspective.
  - Satisfies either [encoding.TextUnmarshaler] or [encoding.BinaryUnmarshaler]
    interfaces, checked for in that order.
  - Is a value type.

When testing, there is no interface for the [Field], so a test struct must be
defined and a call to [SettableFields] is required to get a slice of [Field] to
work with. There are good examples of this in the confhandler package tests.
*/
package stronf
