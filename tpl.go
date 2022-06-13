package olapsql

import (
	"sync"

	"github.com/ahmetb/go-linq/v3"
)

func Parallel(input any, function func(any) any, parallelNumber int) chan any {
	wg := sync.WaitGroup{}

	in := make(chan any)
	out := make(chan any)
	for i := 0; i < parallelNumber; i++ {
		wg.Add(1)
		go func() {
			linq.FromChannel(in).ForEach(func(v any) {
				out <- function(v)
			})
			wg.Done()
		}()
	}

	go linq.From(input).ToChannel(in)
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
