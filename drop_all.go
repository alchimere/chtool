package chtool

import (
	"context"
)

// DropAll consumes all channel elements until close
func DropAll[T any](ctx context.Context, in <-chan (T)) chan (T) {
	return Filter(ctx, in, func(T) bool {
		return false
	})
}
