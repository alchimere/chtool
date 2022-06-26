package chtool

import (
	"context"
)

// Filter selects each element that matches the predicate and forward them to another chan
func Filter[T any](ctx context.Context, in <-chan (T), fn func(T) bool) chan (T) {
	out := make(chan (T))

	go func() {
		defer close(out)

		for {
			select {
			case item := <-in:
				if fn(item) {
					out <- item
				}
			case <-ctx.Done():
				break
			}
		}
	}()
	
	return out
}