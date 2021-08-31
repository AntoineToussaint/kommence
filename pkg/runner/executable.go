package runner

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"io"
	"os/exec"
	"strings"
	"sync"
)

type Executable struct {
	Cmd    string
	Args   []string
	logger *output.Logger
	config *configuration.Executable
}


func NewExecutable(logger *output.Logger, c *configuration.Executable) Runnable {
	args := strings.Split(c.Cmd, " ")
	return &Executable{
		logger: logger,
		config: c,
		Cmd:    args[0],
		Args:   args[1:],
	}
}

func (e *Executable) ID() string {
	return e.config.Name
}

func (e *Executable) Start(ctx context.Context, rec chan output.Message) error {
	cmd := exec.Command(e.Cmd, e.Args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()


	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, _ = io.Copy(output.NewLineBreaker(rec, e.ID()), stdout)
		wg.Done()
	}()
	_, _ = io.Copy(output.NewLineBreaker(rec, e.ID()), stderr)

	wg.Wait()
	return cmd.Wait()
}
