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
	Executables []string
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

const tmpl = `{{if .Timestamp}} [{{.Timestamp}}]{{end}}{{if .Level}} [{{.Level}}]{{end}} {{.Parsed}}`

func (r *Runner) Run(ctx context.Context, cfg Configuration) error {

	var starting []Runnable

	var styler output.Styler
	styles := make(map[string]output.Style)

	for _, executable := range cfg.Executables {
		if c, ok := r.Configuration.Execs.Get(executable); ok {
			exec := NewExecutable(r.Logger, c)
			starting = append(starting, exec)
		}
	}
	maxIDLength := 0
	for _, start := range starting {
		if l := len(start.ID()); l > maxIDLength {
			maxIDLength = l
		}
		styles[start.ID()] = styler.Next()
	}
	padding := PaddedID{Length: maxIDLength}

	go func() {
		for msg := range r.Receiver {
			// Parse message
			parsed := output.ParseToStructured(msg.Content)
			// Render it
			rendered := output.FromTemplate(tmpl, parsed)
			// Style it
			style := styles[msg.ID]
			r.Logger.Printf(padding.ID(msg.ID) + " >" + rendered + "\n", style...)
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

