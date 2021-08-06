package pkg

import (
	"context"
	"fmt"
	"log"
)

type FlowRuntime struct {
	Runs []string
	Kubes []string
}


func Flow(ctx context.Context, c *FlowConfiguration, r FlowRuntime) {
	var sources []Source
	// For formatter, we need the longest run length
	maxLength := 0

	allRunConfigurations := make(map[string]*RunnerConfiguration)
	var runs []string

	for _, run := range c.Runs {
		runs = append(runs, run.Run.Name)
		allRunConfigurations[run.Run.Name] = run
	}
	// If runs are specified, we only start those
	if len(r.Runs) > 0 {
		runs = r.Runs
	}
	for _, run := range runs {
		if len(run) > maxLength {
			maxLength = len(run)
		}
	}


	allKubeConfigurations := make(map[string]*KubeConfiguration)
	var kubes []string

	for _, kube := range c.Kubes {
		kubes = append(kubes, kube.Kube.Name)
		allKubeConfigurations[kube.Kube.Name] = kube
	}
	// If kubes are specified, we only start those
	if len(r.Kubes) > 0 {
		kubes = r.Kubes
	}
	for _, kube := range kubes {
		if len(kube) > maxLength {
			maxLength = len(kube)
		}
	}


	// Load the runners
	for _, run := range runs {
		fmt.Println("running:", run)
		config, ok := allRunConfigurations[run]
		if !ok {
			log.Fatalf("run %v is not in the configuration", run)
		}
		runner, err := NewRunner(*config)
		if err != nil {
			log.Fatal(err)
		}
		formatted, err := NewFormatter(RunType, config.Format, runner, maxLength)
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, formatted)
	}

	// Load the kubes
	for _, kube := range kubes {
		fmt.Println("forwarding:", kube)
		config, ok := allKubeConfigurations[kube]
		if !ok {
			log.Fatalf("kube %v is not in the configuration", kube)
		}
		runner, err := NewKube(*config)
		if err != nil {
			log.Fatal(err)
		}
		formatted, err := NewFormatter(KubeType, config.Format, runner, maxLength)
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

