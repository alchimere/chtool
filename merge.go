package chtool

import (
	"context"
	"sync"
)

// Merge forwards elements of multiple channels into one
func MergeMerge[T any](ctx context.Context, channels ...<-chan (T)) chan (T) {
	out := make(chan (T))

	go func() {
		defer close(out)

		var wg sync.WaitGroup
		for i := 0; i < len(channels); i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				in := channels[idx]
				for {
					select {
					case item, ok := <-in:
						if !ok {
							return
						}
						out <- item
					case <-ctx.Done():
						return
					}
				}
			}(i)
			wg.Wait()
		}
	}()

	return out
}
