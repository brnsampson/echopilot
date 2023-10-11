package option

import (
	"encoding/json"
    "fmt"
)

type Optional[T any] interface {
    Clone() Option[T]
    IsSome() bool
    IsNone() bool
    Unwrap() (*T, error)
    UnsafeUnwrap() T
    Set(T)
    UnwrapOrDefault(T) T
    UnwrapOrElse(func () T) T
}

func NewOption[T any](value T) Option[T] {
    return Option[T]{ inner: &value }
}

func None[T any]() Option[T] {
    return Option[T]{ inner: nil }
}

type Option[T any] struct {
    inner *T
}

// Optional interface
func (o Option[T]) Clone() Option[T] {
    if o.IsNone() {
        return o
    } else {
        // dereference the pointer to force to to allocate new memory
        o.inner = &(*o.inner)
        return o
    }
}

func (o *Option[T]) IsSome() bool {
    if o.inner == nil {
        return false
    }
    return true
}

func (o *Option[T]) IsNone() bool {
    if o.inner == nil {
        return true
    }
    return false
}

func (o *Option[T]) Unwrap() (*T, error) {
    if o.IsSome() {
        i := o.inner
        o.inner = nil
        return i, nil
    }
    return nil, fmt.Errorf("Attempted to unwrap Option with None value")
}

func (o *Option[T]) UnsafeUnwrap() T {
    if o.IsSome() {
        i := o.inner
        o.inner = nil
        return *i
    }
    panic("Attempted to unsafely unwrap an Option with None value")
}

func (o *Option[T]) Set(new T) {
    o.inner = &new
}

func (o *Option[T]) UnwrapOrDefault(def T) T {
    if o.IsSome() {
        return *o.inner
    }
    return def
}

func (o *Option[T]) UnwrapOrElse(f func() T) T {
    if o.IsSome() {
        return *o.inner
    }
    return f()
}

// Unmarshaller interface
func (o *Option[T]) UnmarshalJSON(data []byte) error {
  var tmp T
  if err := json.Unmarshal(data, &tmp); err != nil {
    return err
  }
  *o = NewOption(tmp)
  return nil
}
