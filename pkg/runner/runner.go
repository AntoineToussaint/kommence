package runner

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"sync"
)

type Runner struct {
	Receiver chan output.Message
	Configuration *configuration.Configuration
	Logger        *output.Logger
}

type Configuration struct {
	Runs []string
}

type Runnable interface {
	Start(ctx context.Context, rec chan output.Message) error
	ID() string
}

func New(log *output.Logger, c *configuration.Configuration) Runner {
	return Runner{
		Logger: log,
		Configuration: c,
		Receiver: make(chan output.Message),
	}
}

type PaddedID struct {
	Length int
}

func (p *PaddedID) ID(id string) string {
	padding := ""
	for i := 0; i < p.Length - len(id); i++ {
		padding += " "
	}
	return id + padding
}

func (r *Runner) Run(ctx context.Context, cfg Configuration) error {
	var starting []Runnable

	for _, run := range cfg.Runs {
		if c, ok := r.Configuration.Execs.Get(run); ok {
			exec := NewExecutable(r.Logger, c)
			starting = append(starting, exec)
		}
	}
	maxIDLength := 0
	for _, start := range starting {
		if l := len(start.ID()); l > maxIDLength {
			maxIDLength = l
		}
	}
	padding := PaddedID{Length: maxIDLength}

	go func() {
		for msg := range r.Receiver {
			r.Logger.Printf(padding.ID(msg.ID) + " > " + msg.Content + "\n")
		}
	}()
	var wg sync.WaitGroup
	wg.Add(len(starting))
	for _, start := range starting {
		go func(s Runnable) {
			_ = s.Start(ctx, r.Receiver)
			wg.Done()
		}(start)
	}
	wg.Wait()
	return nil
}

