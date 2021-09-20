package runner

import (
	"context"
	"fmt"
	"github.com/AntoineToussaint/kommence/pkg/configuration"
	"github.com/AntoineToussaint/kommence/pkg/output"
	"github.com/fatih/color"
	"strings"
)

type Runner struct {
	Receiver      chan output.Message
	Configuration *configuration.Configuration
	Logger        *output.Logger
	tasks         []Runnable
}

type Runtime struct {
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

func (r *Runner) Run(ctx context.Context, cfg *Runtime) error {
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
			if msg.Type != output.Log {
				continue
			}
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

	errors := make(chan error)
	for _, task := range r.tasks {
		go func(s Runnable) {
			// Some runnable returns error (Pod) and some don't (Executable)
			// On error, we should return: stop kommence
			err := s.Start(ctx, r.Receiver)
			if err != nil {
				errors <- err
				r.Logger.Printf("%v received an unrecoverable error: %v\n", s.ID(), err)
			}
		}(task)
	}
	// Wait for context Done or if we stop
	for {
		select {
		case <-ctx.Done():
			r.Logger.Debugf("Done with all tasks\n")
			return nil
		case err := <-errors:
			r.Logger.Debugf("Received an error from a task: %v\n", err)
			return err
		}
	}
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
