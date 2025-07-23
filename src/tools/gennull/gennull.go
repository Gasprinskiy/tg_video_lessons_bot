package gennull

type GenericNull[T any] struct {
	Valid bool
	Value T
}

func NewGenericNull[T any](value T) GenericNull[T] {
	return GenericNull[T]{
		Valid: true,
		Value: value,
	}
}
