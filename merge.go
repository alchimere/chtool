package chtool

import (
	"context"
	"sync"
)

// Merge forwards elements of two channels into one
func Merge[T any](ctx context.Context, ch1, ch2 <-chan (T)) chan (T) {
	out := make(chan (T))

	go func() {
		defer close(out)

		for {
			select {
			case item := <-ch1:
				out <- item
			case item := <-ch2:
				out <- item
			case <-ctx.Done():
				break
			}
		}
	}()

	return out
}

// MultiMerge forwards elements of N channels into one
// Using multiple goroutines is simpler than using reflect.Select
func MultiMerge[T any](ctx context.Context, channels ...<-chan (T)) chan (T) {
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
					case item := <-in:
						out <- item
					case <-ctx.Done():
						break
					}
				}
			}(i)
			wg.Wait()
		}
	}()

	return out
}
