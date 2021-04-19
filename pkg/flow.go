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
		formatted, err := NewFormatter(config.Format, runner)
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

