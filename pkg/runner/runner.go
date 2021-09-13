package runner

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"strings"
	"sync"

	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
)

type Runner struct {
	Receiver      chan output.Message
	Configuration *configuration.Configuration
	Logger        *output.Logger
	tasks         []Runnable
}

type Configuration struct {
	Executables []string
	Pods        []string
}

type Runnable interface {
	ID() string
	Start(ctx context.Context, rec chan output.Message) error
	Stop(ctx context.Context, rec chan output.Message) error
}

func New(log *output.Logger, c *configuration.Configuration) *Runner {
	return &Runner{
		Logger:        log,
		Configuration: c,
		Receiver:      make(chan output.Message),
	}
}

type PaddedID struct {
	Length int
}

func (p *PaddedID) ID(id string) string {
	padding := ""
	for i := 0; i < p.Length-len(id); i++ {
		padding += " "
	}
	return id + padding
}

const tmpl = `{{if .Timestamp}} [{{.Timestamp}}]{{end}}{{if .Level}} [{{.Level}}]{{end}} {{.Parsed}}`


func (r *Runner) Run(ctx context.Context, cfg *Configuration) error {
	var styler output.Styler
	styles := make(map[string]output.Style)

	for _, executable := range cfg.Executables {
		if c, ok := r.Configuration.Execs.Get(executable); ok {
			exec := NewExecutable(r.Logger, c)
			r.tasks = append(r.tasks, exec)
		}
	}

	// Load Kubernetes client
	if len(cfg.Pods) > 0 {
		r.Logger.Debugf("loading kubernetes client\n")
		LoadKubeClient()
	}

	for _, pod := range cfg.Pods {
		if c, ok := r.Configuration.Pods.Get(pod); ok {
			exec := NewPod(r.Logger, c)
			r.tasks = append(r.tasks, exec)
		}
	}

	if len(r.tasks) == 0 {
		r.Logger.Printf("Nothing to run/forward\n", color.Bold)
		return nil
	}

	// Figure out padding and styles
	maxIDLength := 0
	for _, start := range r.tasks {
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
			rendered := output.FromTemplate(r.Logger, tmpl, parsed)
			// Style it
			style := styles[msg.ID]
			// Regular message
			r.Logger.Printf(padding.ID(msg.ID)+" >"+rendered+"\n", style...)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(len(r.tasks))

	for _, task := range r.tasks {
		go func(s Runnable) {
			err := s.Start(ctx, r.Receiver)
			if err != nil {
				r.Logger.Errorf(err.Error() + "\n")
			}
			wg.Done()
		}(task)
	}
	wg.Wait()
	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	var errors []string
	for _, task := range r.tasks {
		err := task.Stop(ctx, r.Receiver)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("can't stop properly: %v", strings.Join(errors, ", "))
	}
	return nil
}
