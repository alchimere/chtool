package chtool

import (
	"context"
)

// Each applies func on each elements and forward all these elements to out channel
func Each[T any](ctx context.Context, in <-chan (T), fn func(T)) chan (T) {
	return Map(ctx, in, func(item T) T {
		fn(item)
		return item
	})
}
