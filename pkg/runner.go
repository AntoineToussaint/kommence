package pkg

import (
	"context"
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os/exec"
	"strings"
	"time"
)

type RunnerConfiguration struct {
	Name  string
	Command   string
	Watch []string
}

type Runner struct {
	RunnerConfiguration
}

func (r Runner) ID() string {
	return r.Name
}

func (r Runner) run(out chan string) {
	args := strings.Split(r.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = NewLineBreaker(out)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR", err)
	}
}


func watch(config RunnerConfiguration) <-chan bool {
	w := watcher.New()
	w.SetMaxEvents(1)
	events := make(chan bool)
	//r := regexp.MustCompile(".*.go")
	//w.AddFilterHook(watcher.RegexFilterHook(r, false))
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event)
				events <- true
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	for _, path := range config.Watch {
		if err := w.AddRecursive(path); err != nil {
			log.Fatalln(err)
		}
	}
	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}

	}()
	return events
}

func (r Runner) Produce(ctx context.Context) <-chan string {
	out := make(chan string)
	events := watch(r.RunnerConfiguration)
	go func() {
		r.run(out)
		for {
			select {
			case <-events:
				go r.run(out)
			}
		}
	}()
	return out
}
