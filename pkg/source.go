package pkg

import (
	"context"
	"fmt"
	"sync"
)

type Source interface {
	ID() string
	Produce(ctx context.Context) <-chan string
}


type Decorate struct {
	source Source
}

func (d Decorate) Produce(ctx context.Context) <-chan string {
	out := make(chan string)
	go func() {
		for output := range d.source.Produce(ctx) {
			out <- fmt.Sprintf("(%v) %v", d.source.ID(), output)
		}
	}()
	return out
}



func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

