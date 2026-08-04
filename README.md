# `go-exception` — Behavior Summary

Module `github.com/thanhminhmr/go-exception`, Go 1.25, MPL-2.0. Optional zerolog integration (disabled with the
`no_zerolog` build tag).

## Core contract

`Exception` is a **sealed interface** (`exception.go:32`): a private `__()` method prevents any external implementation.
All methods use **value receivers**; mutation methods **return a new `Exception`** rather than mutating in place.
Documented rule: **always use the returned value** — never assume the original is unchanged.

## The three implementations

The interface has three concrete types, chosen by weight (cheapest sufficient form wins):

### `String` (`string.go:38`) — lightweight, constant-friendly

`type String string`. Carries only type + message, parsed by the `": "` separator (`string.go:15`):

- Missing separator ⇒ whole string is the type, message is empty.
- `Error()` drops the separator when type or message is empty (e.g. `"IOError: msg"` → `"IOError: msg"`; `": msg"` →
  `"msg"`; `"IOError:"` → `"IOError"`).
- Usable as a `const` (e.g. `const ErrRead = exception.String("IOError: read failed")`).
- All field getters (`GetCause`, `GetSuppressed`, `GetRecovered`, `GetStackTrace`, `GetExtras`, `GetExtra`) return zero
  values; `Clone()` returns itself (immutable).

### `multipleErrors` (`multiple_errors.go:37`) — the join form

`type multipleErrors []error`. Empty type and message; `GetCause()` and `Unwrap() []error` both return itself. `Error()`
is `fmt.Sprintf("%v", []error(e))`. Produced by `Join`.

### `fullException` (`full_exception.go:18`) — the full form

Struct with `Type`, `Message`, `Cause` (`multipleErrors`), `Suppressed` (`multipleErrors`), `Recovered any`,
`StackTrace StackFrames`, `Extras map[string]any`. `Error()` returns `"Type: Message"`, or whichever of type/message is
non-empty, or `""` if both empty. `Unwrap() []error` returns `Cause`.

## Promotion-on-enrichment

`String` and `multipleErrors` **stay cheap until enriched**, then promote to `fullException` carrying over their
existing data:

| Op on `String`              | Op on `multipleErrors`      | Result                                                                                             |
|-----------------------------|-----------------------------|----------------------------------------------------------------------------------------------------|
| `AddCause` (non-empty)      | `AddCause` (non-empty)      | `fullException` with `Cause` set, type preserved (`String`) or causes preserved (`multipleErrors`) |
| `AddSuppressed` (non-empty) | `AddSuppressed` (non-empty) | `fullException` with `Suppressed` set                                                              |
| `SetRecovered(non-nil)`     | `SetRecovered(non-nil)`     | `fullException` with `Recovered` set                                                               |
| `FillStackTrace`            | `FillStackTrace`            | `fullException` with `StackTrace` set                                                              |
| `SetExtras(non-nil)`        | `SetExtras(non-nil)`        | `fullException` with `Extras` set                                                                  |
| `SetExtra(key, non-nil)`    | `SetExtra(key, non-nil)`    | `fullException` with a fresh `Extras` map containing only that key                                 |
| `SetMessage` (non-empty)    | `SetMessage` (non-empty)    | `String` returns a new `String`; `multipleErrors` promotes to `fullException` keeping `Cause`      |

Once a `fullException`, all further mutations stay on `fullException` (struct is copied and the field is reassigned).

## Construction helpers

### `Join(errors ...error) Exception` (`multiple_errors.go:25`)

Combines errors into a `multipleErrors`. **Nil inputs are filtered.** Nested `multipleErrors` are flattened
(auto-unboxed). Returns `nil` if nothing remains.

### `Template` (`template.go`)

`type Template string`. `Format(args...)` returns `String(fmt.Sprintf(t, args...))`. Reusable constant format for
exception messages.

### `Panic` / `Recover` (`panic.go`)

- `PanicError = String("panicked")` — the type used for panic-originated exceptions.
- `Panic(v)`: if `v` is already an `Exception` whose type is `PanicError`, re-panics it unchanged (enables chained
  recover handlers without altering state). Otherwise wraps `v` as
  `fullException{Type: "panicked", Recovered: v, StackTrace: StackTrace(1)}` and panics it.
- `Recover(callback)`: no-op when no panic occurred. If the recovered value is an `Exception` with type `PanicError`,
  passes it to the callback unchanged. Otherwise builds a `fullException` and **strips leading `runtime.*panic*`
  frames** from the trace (`panic.go:87-97`) so the trace starts at the real panic site. `Recover(nil)` panics with
  `"BUG: callback is nil"`. Intended for `defer exception.Recover(func(ex exception.Exception) { ... })`.

## Field semantics

- **Cause** (`GetCause`/`AddCause`): the root errors that led to this exception. On `fullException`, `AddCause` appends
  via `multipleErrors.append` (`multiple_errors.go:156`) which filters nils and flattens nested `multipleErrors`.
- **Suppressed** (`GetSuppressed`/`AddSuppressed`): errors intentionally ignored or deferred while handling this
  exception. Same append semantics as Cause.
- **Recovered** (`GetRecovered`/`SetRecovered`): the value captured from a panic.
- **StackTrace** (`GetStackTrace`/`FillStackTrace`): `StackFrame{Function, File, Line}` slice. `FillStackTrace(skip)`
  calls `StackTrace(skip+1)` so it captures from the caller of `FillStackTrace`; `skip=0` includes that caller.
  `StackTrace` captures up to 64 PCs via `runtime.Callers(2+skip, …)` and appends **every** frame returned by
  `runtime.CallersFrames.Next`, including the last (where `more==false`).
- **Extras** (`GetExtras`/`SetExtras`/`GetExtra`/`SetExtra`): arbitrary key-value metadata. `SetExtra` lazily allocates
  the map on first write; `SetExtra(key, nil)` deletes the key (no-op on a nil map, since `delete` on nil is allowed in
  Go).

## `errors` package integration

- `errors.Is` walks the cause chain via multi-`Unwrap() []error`: `fullException.Unwrap()` returns `Cause`,
  `multipleErrors.Unwrap()` returns itself. A `Join` result therefore matches any of its members.
- `errors.As` works through the same Unwrap chain.

## `Clone`

- `String.Clone()` → returns itself (immutable).
- `multipleErrors.Clone()` → `slices.Clone` (new slice, shared error elements).
- `fullException.Clone()` → new struct with `slices.Clone` for Cause/Suppressed/StackTrace and `maps.Clone` for Extras.
  **Referenced values (error elements, the recovered value) are not deeply cloned** — they're shared between original
  and clone.

## Nil handling invariants

- `AddCause`/`AddSuppressed`/`Join` silently filter all nil error arguments.
- On `String`/`multipleErrors`: `SetRecovered(nil)`, `SetExtras(nil)`, and `SetExtra(key, nil)` are **no-ops** returning
  the receiver unchanged (no promotion).
- `GetExtra` on an exception with no extras returns `(nil, false)`.

## Zerolog integration (optional)

`zerolog.go` compiles unless the `no_zerolog` build tag is set. `String`, `fullException`, `multipleErrors`,
`StackFrame`, and `StackFrames` all implement `MarshalZerologObject` / `MarshalZerologArray`, omitting empty fields and
choosing `AnErr` (single) vs `Errs` (multiple) based on slice length.
