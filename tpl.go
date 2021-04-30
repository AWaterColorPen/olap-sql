package olapsql

import (
	"sync"

	"github.com/ahmetb/go-linq/v3"
)

func Parallel(input interface{}, function func(interface{}) interface{}, parallelNumber int) chan interface{} {
	wg := sync.WaitGroup{}

	in := make(chan interface{})
	out := make(chan interface{})
	for i := 0; i < parallelNumber; i++ {
		wg.Add(1)
		go func() {
			linq.FromChannel(in).ForEach(func(v interface{}) {
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
