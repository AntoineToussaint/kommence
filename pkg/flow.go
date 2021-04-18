package pkg

import (
	"context"
	"fmt"
)

type FlowRuntime struct {
	Runs []string
}

func Flow(c FlowConfiguration, r FlowRuntime) {
	var sources []Source
	allRunConfigurations := make(map[string]RunnerConfiguration)
	for _, run := range c.Run {
		allRunConfigurations[run.Name] = run
	}
	for _, run := range r.Runs {
		sources = append(sources, Runner{RunnerConfiguration: allRunConfigurations[run]})
	}
	fmt.Printf("running flow with %v sources\n", len(sources))
	out := Console{}
	ctx := context.Background()
	out.Consume(ctx, sources)
}
