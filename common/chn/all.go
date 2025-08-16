package chn

func All[T any](ch <-chan T) (xs []T) {
	for x := range ch {
		xs = append(xs, x)
	}
	return xs
}
