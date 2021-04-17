package pkg

import (
	"context"
	"fmt"
)

type Sink interface {
	Consume(ctx context.Context, sources []Source)
}

type Console struct {}

func (c *Console) Consume(ctx context.Context, sources []Source) {
	var all []<-chan string
	for _, source := range sources {
		all = append(all, Decorate{source}.Produce(ctx))
	}
	out := merge(all...)
	for output := range out {
		fmt.Println(output)

	}
}