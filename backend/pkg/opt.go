package pkg

type Opt[T any] struct {
	value  T
	exists bool
}

func NewEmpty[T any]() Opt[T] {
	return Opt[T]{
		exists: false,
	}
}

func NewWithValue[T any](value T) Opt[T] {
	return Opt[T]{
		value:  value,
		exists: true,
	}
}

func (o *Opt[T]) Get() (T, bool) {
	return o.value, o.exists
}

func (o *Opt[T]) Set(value T) {
	o.value = value
	o.exists = true
}

func (o *Opt[T]) Unset() {
	o.exists = false
}

func (o *Opt[T]) IsExists() bool {
	return o.exists
}
