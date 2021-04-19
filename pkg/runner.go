package pkg

import (
	"context"
	"fmt"
	"github.com/radovskyb/watcher"
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
	Format FormatterConfiguration
}

type Runner struct {
	RunnerConfiguration
	cmd *exec.Cmd
	out chan Message
}

func NewRunner(c RunnerConfiguration) (*Runner, error) {
	out := make(chan Message)
	runner := Runner{RunnerConfiguration: c, out: out}
	return &runner, nil
}

func (r *Runner) ID() string {
	return r.Name
}

func (r *Runner) Produce(ctx context.Context) <-chan Message {
	return r.out
}

func (r *Runner) CreateWatcher() <-chan bool {
	out := make(chan bool)
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify rename and move events.
	w.FilterOps(watcher.Write)

	// Only files that match the regular expression during file listings
	// will be watched.
	//r := regexp.MustCompile("^abc$")
	//w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case <-w.Event:
				out <- true
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	for _, p := range r.Watch {
		if err := w.AddRecursive(p); err != nil {
			log.Fatalln(err)
		}
	}

	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
	return out
}

func (r *Runner) Start(ctx context.Context) {
	go r.Run(ctx)
	// Watcher
	if w := r.CreateWatcher(); w != nil {
		for {
			select {
			case <-w:
				r.Restart(ctx)
			}
		}
	}

}


func (r *Runner) Run(ctx context.Context) {
	args := strings.Split(r.Command, " ")
	r.cmd = exec.Command(args[0], args[1:]...)
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
	fmt.Println("restarting:", r.Name)
	go r.Start(ctx)
}

