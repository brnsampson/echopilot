package option

import (
	"encoding/json"
	"fmt"
)

type NoneError struct {
	msg string
}

func NewNoneError(msg string) *NoneError {
	return &NoneError{msg}
}

func (e NoneError) Error() string {
	return e.msg
}

type primatives interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~bool | ~string
}

type Optional[T comparable] interface {
	IsSome() bool
	IsNone() bool
	Clear()
	Set(T)
	Unwrap() (T, error)
	UnsafeUnwrap() T
	UnwrapOrDefault(T) T
	UnwrapOrElse(func() T) T
	Match(T) bool
	Eq(Optional[T]) bool
	Transform(func(T) (T, error)) error
	TransformOr(func(T) (T, error), T)
	TransformOrElse(func(T) (T, error), func() T)
	BinaryTransform(second T, f func(T, T) error) error
}

type ConfigOptional[T primatives] interface {
	Optional[T]

	// Satisfies encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
}

type Option[T comparable] struct {
	inner T
	none  bool
}

func Some[T comparable](value T) Option[T] {
	return Option[T]{inner: value, none: false}
}

func None[T comparable]() Option[T] {
	var tmp T
	return Option[T]{inner: tmp, none: true}
}

func (o Option[T]) IsSome() bool {
	return !o.none
}

func (o Option[T]) IsNone() bool {
	return o.none
}

func (o *Option[T]) Clear() {
	o.none = true
}

func (o *Option[T]) Set(value T) {
	o.inner = value
	o.none = false
}

func (o *Option[T]) Unwrap() (T, error) {
	if o.IsSome() {
		o.none = true
		return o.inner, nil
	}
	return o.inner, fmt.Errorf("Attempted to unwrap Option with None value")
}

func (o *Option[T]) UnsafeUnwrap() T {
	if o.IsSome() {
		o.none = true
		return o.inner
	}
	panic("Attempted to unsafely unwrap an Option with None value")
}

func (o *Option[T]) UnwrapOrDefault(def T) T {
	if o.IsSome() {
		return o.inner
	}
	return def
}

func (o *Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsSome() {
		return o.inner
	}
	return f()
}

func (o Option[T]) Match(probe T) bool {
	if o.none {
		return false
	} else {
		return o.inner == probe
	}
}

func (o Option[T]) Eq(other Optional[T]) bool {
	if o.none && other.IsNone() {
		return true
	} else if !o.none && other.IsSome() {
		// We do not know if other is a pointer type or not, so play it safe
		return other.Match(o.inner)
	} else {
		// one is none and the other is some
		return false
	}
}

func (o *Option[T]) Transform(f func(T) (T, error)) error {
	if o.IsNone() {
		return NewNoneError("Attempted to transform None value")
	}

	tmp, err := f(o.inner)
	if err != nil {
		return err
	}

	o.inner = tmp
	return nil
}

func (o *Option[T]) TransformOr(f func(T) (T, error), backup T) {
	if o.IsNone() {
		o.inner = backup
	} else {
		tmp, err := f(o.inner)
		if err != nil {
			o.inner = backup
		} else {
			o.inner = tmp
		}
	}
}

func (o *Option[T]) TransformOrElse(f func(T) (T, error), backup func() T) {
	if o.IsNone() {
		o.inner = backup()
	} else {
		tmp, err := f(o.inner)
		if err != nil {
			o.inner = backup()
		} else {
			o.inner = tmp
		}
	}
}

func (o *Option[T]) BinaryTransform(second T, f func(T, T) (T, error)) error {
	if o.IsNone() {
		return NewNoneError("Attempted to transform None value")
	}

	tmp, err := f(o.inner, second)
	if err != nil {
		return err
	}

	o.inner = tmp
	return nil
}

// Unmarshaller interface
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var tmp T
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	o.Set(tmp)
	return nil
}
