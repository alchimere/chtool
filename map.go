package chtool

import (
	"context"
)

// Map applies tranformation on each channel element and produce a channel with result
func Map[T, U any](ctx context.Context, in <-chan (T), fn func(T) U) chan (U) {
	out := make(chan (U))

	go func() {
		defer close(out)

		for {
			select {
			case item, ok := <-in:
				if !ok {
					return
				}
				out <- fn(item)
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
