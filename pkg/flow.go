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
	allRunConfigurations := make(map[string]RunnerConfiguration)
	for _, run := range c.Run {
		allRunConfigurations[run.Name] = run
	}
	for _, run := range r.Runs {
		config, ok := allRunConfigurations[run]
		if !ok {
			log.Fatalf("run %v is not in the configuration", run)
		}
		runner, err := NewRunner(config)
		if err != nil {
			log.Fatal(err)
		}
		go runner.Do(ctx)
		sources = append(sources, runner)
	}
	fmt.Printf("running flow with %v sources\n", len(sources))
	out := Console{}
	out.Consume(ctx, sources)
}

