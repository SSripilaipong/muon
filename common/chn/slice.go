package chn

func FromSlice[T any](xs []T) <-chan T {
	ch := make(chan T, len(xs))
	defer close(ch)
	for _, x := range xs {
		ch <- x
	}
	return ch
}
