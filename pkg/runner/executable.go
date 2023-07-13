package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/AntoineToussaint/kommence/pkg/configuration"
	"github.com/AntoineToussaint/kommence/pkg/output"
	"github.com/radovskyb/watcher"
)

type Executable struct {
	cmd        string
	args       []string
	command    *exec.Cmd
	stdErrMode string

	logger      *output.Logger
	config      *configuration.Executable
	restartChan chan interface{}
}

func NewExecutable(logger *output.Logger, c *configuration.Executable) Runnable {
	args := strings.Split(c.Cmd, " ")
	return &Executable{
		logger:     logger,
		config:     c,
		cmd:        args[0],
		args:       args[1:],
		stdErrMode: c.StdErr,
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
		e.logger.Debugf("start: %v\n", e.ID())
		e.start(ctx, rec)
	}()
	for {
		select {
		case <-w.Event:
			e.logger.Debugf("watcher caused restart: %v\n", e.ID())
			go e.restart(ctx, rec)
		case err := <-w.Error:
			e.logger.Errorf("watcher error: %v: %v\n", e.ID(), err)
			return err
		case <-w.Closed:
			e.logger.Debugf("watcher closed: %v\n", e.ID())
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}

func (e *Executable) Stop(ctx context.Context, rec chan output.Message) error {
	e.logger.Debugf("stopping: %v\n", e.ID())
	return e.kill(ctx, rec)
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
	if e.config.Delay != "" {
		e.logger.Debugf("delaying %v %v\n", e.ID(), e.config.Delay)
		d, _ := time.ParseDuration(e.config.Delay)
		time.Sleep(d)
	}
	e.logger.Debugf("starting %v\n", e.ID())
	e.command = exec.CommandContext(ctx, e.cmd, e.args...)
	e.command.Env = os.Environ()
	for k, v := range e.config.Env {
		e.command.Env = append(e.command.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Request the OS to assign process group to the new process, to which all its children will belong
	e.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := e.command.StdoutPipe()
	stderr, _ := e.command.StderrPipe()

	if err := e.command.Start(); err != nil {
		e.logger.Errorf("can't start %v: %v", e.ID(), err)
	}
	// Export resources
	go func() {
		// TODO Doesn't work on MAC
		//e.MonitorMemory(e.command.Process.Pid, rec)
	}()

	// Export logs
	go func() {
		_, _ = io.Copy(output.NewLineBreaker(rec, e.ID(), output.Log), stdout)
	}()
	if e.stdErrMode == configuration.AsLog {
		_, _ = io.Copy(output.NewLineBreaker(rec, e.ID(), output.Log), stderr)
	}
	return
}

func (e *Executable) kill(ctx context.Context, rec chan output.Message) error {
	rec <- output.Message{ID: e.ID(), Type: output.Stop, Content: "Stopping"}
	if e.command == nil || e.command.Process == nil {
		return nil
	}
	if err := syscall.Kill(-e.command.Process.Pid, syscall.SIGKILL); err != nil {
		e.logger.Errorf("failed to kill process %v: %v\n", e.ID(), err)
		return err
	}
	return nil
}

func (e *Executable) MonitorMemory(pid int, rec chan output.Message) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			mem, err := calculateMemory(pid)
			if err != nil {
				e.logger.Errorf("can't calc memory for %v: %v", e.ID(), err)
				continue
			}
			rec <- output.Message{
				ID:      e.ID(),
				Type:    output.Memory,
				Content: fmt.Sprintf("%v", mem),
			}
		}
	}
}

func calculateMemory(pid int) (uint64, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := uint64(0)
	pfx := []byte("Pss:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			res += size
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}
	return res, nil
}

func (e *Executable) restart(ctx context.Context, rec chan output.Message) {
	err := e.kill(ctx, rec)
	if err != nil {
		e.logger.Errorf("can't kill %v: %v\n", e.ID(), err)
	}
	rec <- output.Message{ID: e.ID(), Type: output.Restart, Content: "⏳ RESTARTING ⏳"}
	e.start(ctx, rec)
}
