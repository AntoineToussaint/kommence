package pkg

import (
	"context"
	"fmt"
	"log"
)

type FlowRuntime struct {
	Runs []string
}


func Flow(ctx context.Context, c *FlowConfiguration, r FlowRuntime) {
	var sources []Source
	allRunConfigurations := make(map[string]*RunnerConfiguration)
	var runs []string

	for _, run := range c.Run {
		runs = append(runs, run.Run.Name)
		allRunConfigurations[run.Run.Name] = run
	}
	// If runs are specified, we only start those
	if len(r.Runs) > 0 {
		runs = r.Runs
	}
	// For formatter, we need the longest run length
	maxLength := 0
	for _, run := range runs {
		if len(run) > maxLength {
			maxLength = len(run)
		}
	}
	// Load the runners
	for _, run := range runs {
		config, ok := allRunConfigurations[run]
		if !ok {
			log.Fatalf("run %v is not in the configuration", run)
		}
		runner, err := NewRunner(*config)
		if err != nil {
			log.Fatal(err)
		}
		formatted, err := NewFormatter(config.Run.Format, runner, maxLength)
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, formatted)
	}
	fmt.Printf("running flow with %v sources\n", len(sources))

	for _, source := range sources {
		go source.Start(ctx)
	}
	out := Console{}
	out.Consume(ctx, sources)
}

