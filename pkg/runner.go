package pkg

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)


/*
A Runner is an example of a Source

 */
type RunnerConfiguration struct {
	Name    string
	Command string
	Watch   []string
}

type Runner struct {
	RunnerConfiguration
	cmd *exec.Cmd
	out chan Message
}

func NewRunner(c RunnerConfiguration) (*Runner, error) {
	out := make(chan Message)
	args := strings.Split(c.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	runner := Runner{RunnerConfiguration: c, out: out, cmd: cmd}
	return &runner, nil
}

func (r *Runner) ID() string {
	return r.Name
}

func (r *Runner) Produce(ctx context.Context) <-chan Message {
	return r.out
}

func (r *Runner) Do(ctx context.Context) {
	go r.Start(ctx)
	t := time.NewTimer(5 * time.Second)
	<- t.C
	if strings.Contains(r.Name,"counter") {
		fmt.Println("restarting")
		r.Restart(ctx)
	}

}


func (r *Runner) Start(ctx context.Context) {
	stdout, _ := r.cmd.StdoutPipe()
	stderr, _ := r.cmd.StderrPipe()
	err := r.cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, err = io.Copy(NewLineBreaker(r.out), stdout)
		wg.Done()
	}()
	_, _ = io.Copy(NewLineBreaker(r.out), stderr)

	wg.Wait()
	_ = r.cmd.Wait()
}

func (r *Runner) Restart(ctx context.Context) {
	if err := r.cmd.Process.Kill(); err != nil {
		log.Fatalf("failed to kill process %v: %v", r.Name, err)
	}
	go r.Start(ctx)
}

