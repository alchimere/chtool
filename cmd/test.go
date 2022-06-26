package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alchimere/chtool"
)

// func Filter[T any](in <-chan (T), fn func(T) bool) chan (T) {
// 	out := make(chan T)
// 	go func() {
// 		defer close(out)

// 		var wg sync.WaitGroup
// 		for i := 0; i < 2; i++ {
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				for item := range in {
// 					if fn(item) {
// 						out <- item
// 					}
// 				}
// 			}()
// 		}

// 		wg.Wait()
// 		fmt.Printf("Filter done\n")
// 	}()
// 	return out
// }

// func Map[T, U any](in <-chan (T), fn func(T) U) chan (U) {
// 	out := make(chan (U))
// 	go func() {
// 		defer close(out)

// 		var wg sync.WaitGroup
// 		for i := 0; i < 2; i++ {
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				for item := range in {
// 					out <- fn(item)
// 				}
// 			}()
// 		}

// 		wg.Wait()
// 		fmt.Printf("Map done\n")
// 	}()
// 	return out
// }

func ToChannel[T any](s <-chan (T), ch chan<- (T)) {
	defer close(ch)
	for item := range s {
		ch <- item
	}
}

func FromSlice[T any](slice []T) chan (T) {
	ch := make(chan (T))
	go func() {
		defer close(ch)
		for _, t := range slice {
			ch <- t
		}
		fmt.Printf("Input stream done\n")
	}()
	return ch
}

func functionalStyle(out chan<- (string)) {
	ctx := context.TODO()
	ch := FromSlice([]int{1, 2, 3, 4, 5, 6})

	chStr := chtool.Map(ctx, ch,
		func(n int) string {
			//<-time.After(time.Duration(1000-n*100) * time.Millisecond)
			return fmt.Sprintf("page %d\n", n)
		},
		//chtool.Parallel(2),
		//chtool.PreserveOrder()
	)

	chStr2 := chtool.Filter(ctx, chStr,
		func(s string) bool {
			//<-time.After(time.Duration(1000-n*100) * time.Millisecond)
			return !strings.Contains(s, "5")
		},
		//chtool.Parallel(2),
		//chtool.PreserveOrder()
	)

	chDone := chtool.Each(ctx, chStr2,
		func(str string) {
			fmt.Printf("done %s\n", str)
		},
	)

	ToChannel(chDone, out)
}

func main() {
	ch := make(chan (string))
	go func() {
		for s := range ch {
			fmt.Printf("main: %s\n", s)
		}
	}()

	functionalStyle(ch)

	<-time.After(5 * time.Second)
}
