package main

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg"
)



func main() {
	runner := pkg.Runner{RunnerConfiguration: pkg.RunnerConfiguration{Name: "ls", Cmd: "ls .", Watch: []string{"pkg"}}}
	console := pkg.Console{}
	ctx := context.Background()
	console.Consume(ctx, []pkg.Source{runner})
}
