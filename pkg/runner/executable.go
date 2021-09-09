package runner

import (
	"context"
	"fmt"
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/radovskyb/watcher"
	"io"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Executable struct {
	cmd     string
	args    []string
	command *exec.Cmd

	logger  *output.Logger
	config  *configuration.Executable
}

func NewExecutable(logger *output.Logger, c *configuration.Executable) Runnable {
	args := strings.Split(c.Cmd, " ")
	return &Executable{
		logger: logger,
		config: c,
		cmd:    args[0],
		args:   args[1:],
	}
}

func (e *Executable) ID() string {
	return fmt.Sprintf("⚙️ %v", e.config.ID)
}

func (e *Executable) Start(ctx context.Context, rec chan output.Message) error {
	// Watcher
	e.logger.Debugf("creating watcher: %v\n", e.ID())
	w := e.createWatcher()
	go func() {
		for {
			select {
			case <-w.Event:
				go e.restart(ctx, rec)
			case <-ctx.Done():
				e.logger.Debugf("killing process: %v\n", e.ID())
				e.kill(ctx, rec)
				return
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				e.logger.Debugf("watcher closed: %v", e.ID())
				return
			}
		}
	}()
	e.start(ctx, rec)
	return nil
}

func (e *Executable) createWatcher() *watcher.Watcher {
	//out := make(chan bool)
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
	//w.AddFilterHook(watcher.Executables.RegexFilterHook(r, false))

	// Watch this folder for changes.
	for _, p := range e.config.Watch {
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
	return w
}

func (e *Executable) start(ctx context.Context, rec chan output.Message) {
	e.command = exec.CommandContext(ctx, e.cmd, e.args...)
	
	// Request the OS to assign process group to the new process, to which all its children will belong
	e.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := e.command.StdoutPipe()
	stderr, _ := e.command.StderrPipe()

	if err := e.command.Start(); err != nil {
		return
	}
	go func() {
		_, _ = io.Copy(output.NewLineBreaker(rec, e.ID(), false), stdout)
	}()
	_, _ = io.Copy(output.NewLineBreaker(rec, e.ID(), true), stderr)
	e.logger.Debugf("done starting: %v\n", e.ID())
}


func (e *Executable) kill(ctx context.Context, rec chan output.Message) {
	if err := syscall.Kill(-e.command.Process.Pid, syscall.SIGKILL); err != nil {
		e.logger.Errorf("failed to kill process %v: %v\n", e.ID(), err)
	}
}

func (e *Executable) restart(ctx context.Context, rec chan output.Message) {
	e.kill(ctx, rec)
	rec <- output.Message{ID: e.ID(), Content: "** restarting **"}
	e.start(ctx, rec)
}
