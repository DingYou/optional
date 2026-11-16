# optional

A generic optional value for Go 1.18 and later.

```sh
go get github.com/dingyou/optional
```

## Values and presence

```go
import "github.com/dingyou/optional"

var empty optional.Optional[int] // The zero value is empty.
value, ok := empty.Get()         // 0, false

value, ok = optional.Of(0).Get()          // 0, true
flag, present := optional.Of(false).Get() // false, true

values := map[string]*int{"present nil": nil}
pointer, found := values["present nil"]
result := optional.FromOK(pointer, found) // Present, even though pointer is nil.
_, present = result.Get()                // nil, true
```

- `Empty[T]()` and a zero-value `Optional[T]` are empty.
- `Of(value)` is present, including `0`, `false`, and `""`. It panics with
  `ErrNil` for nil values, including typed nil pointers, slices, maps, channels,
  and functions.
- `OfNullable(value)` is empty for nil or typed nil values, and present otherwise.
- `FromOK(value, ok)` trusts `ok`, not nil detection. A true `ok` makes even a
  nil or typed nil value present. A false `ok` discards the value and returns
  `Empty[T]()`.
- `Get()` returns `(value, true)` when present. For every empty Optional, it
  returns `(zero value of T, false)`.

`Val()` retains its existing behavior: it returns the stored value without
checking presence. In particular, `OfNullable[any](typedNil).Val()` retains the
boxed typed nil, whereas `Get()` returns a nil interface and false. `MustVal()`
returns the stored value when present and panics with `ErrNil` when empty.

```go
var pointer *int
nullable := optional.OfNullable[any](pointer)
stored := nullable.Val()         // Retains the typed nil: stored != nil.
value, ok := nullable.Get()       // nil, false
present := optional.FromOK[any](pointer, true)
value, ok = present.Get()         // typed nil, true
```

Use `IsPresent()` / `IsEmpty()` to test presence, `OrElse()` / `OrElseGet()` for
fallbacks, `IfPresent()` for an action, and `Filter()` to retain matching values.
Empty values skip callbacks, except for the fallback supplier in `OrElseGet()`.

## Map and FlatMap

Both functions can change the value's type and skip their callback when the
input is empty.

```go
text := optional.Map(optional.Of(42), func(n int) string {
    return strconv.Itoa(n)
}) // Optional[string] containing "42"; import "strconv" for this example.

flag := optional.FlatMap(text, func(s string) optional.Optional[bool] {
    if s == "42" {
        return optional.Of(false) // Present false, not empty.
    }
    return optional.Empty[bool]()
}) // Optional[bool], never Optional[Optional[bool]].
```

`Map` treats a nil callback result as empty, as if using `OfNullable`.
`FlatMap` returns the callback's Optional unchanged, preserving both its value
and presence, including a nil value made present with `FromOK`.

### Breaking change: FlatMap callback

The main `FlatMap` now accepts `func(T) Optional[U]` instead of
`func(Optional[T]) U`. Update callbacks to accept the contained value and return
an Optional; the result is returned directly, without another Optional wrapper.
Code relying on a zero-value `Optional[T]` being present must also initialize
that value explicitly, for example with `Of(0)`.

Only the main `FlatMap` signature is changed. Existing `MapWithCtx`,
`FlatMapWithCtx`, `MapWithErr`, `FlatMapWithErr`, `MapWithCtxAndErr`, and
`FlatMapWithCtxAndErr` signatures and callback behavior are preserved. In
particular, the `FlatMapWith*` variants still accept an `Optional[T]` callback
argument and wrap its result using `OfNullable`; they do not have the new main
`FlatMap` flattening behavior. All variants now correctly skip zero-value empty
inputs.

## Checks

```sh
go test ./...
go test -race ./...
go vet ./...
```
