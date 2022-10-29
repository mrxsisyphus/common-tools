package map_helper

// Pair 对 表示一组键值对
type Pair[T any, R any] struct {
	First  T
	Second R
}

func NewPair[T any, R any](t T, r R) *Pair[T, R] {
	return &Pair[T, R]{
		First:  t,
		Second: r,
	}
}
