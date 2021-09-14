package cmd

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/antoinetoussaint/kommence/pkg/runner"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var interactiveExecs, interactivePods, interactiveFlows bool
var execs, pods, flows []string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Executables, Pod forwards or Flows",
	Run: func(cmd *cobra.Command, args []string) {
		cancel := make(chan os.Signal, 2)
		signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
		ctx := context.Background()
		ctx, stop := context.WithCancel(ctx)

		log := output.NewLogger(debug)

		log.Debugf("starting in debug mode\n")
		config, err := configuration.Load(log, kommenceDir)
		if err != nil {
			log.Errorf(err.Error(), color.FgRed, color.Bold)
			os.Exit(1)
		}

		var r *runner.Runner
		var c *runner.Runtime

		if interactiveExecs || interactivePods || interactiveFlows {
			log.Debugf("starting interactive mode\n")
			r, c = startInteractive(ctx, log, config)
		} else if len(execs) > 0 || len(pods) > 0 || len(flows) > 0 {
			log.Debugf("starting command line mode\n")
			r, c = startCommandLine(ctx, log, config)
		} else {
			log.Printf("Please specify executables, pods or flows or run in interactive mode.\n")
			os.Exit(0)
		}
		log.Debugf("using runner configuration: %v", c)
		go func() {
			log.Debugf("starting runner\n")
			r.Run(ctx, c)
			// Stop when we are done
			stop()
		}()

		for {
			select {
			case <-cancel:
				log.Printf("Stopping kommence.\n", color.Bold)
				r.Stop(ctx)
				os.Exit(0)

			case <-ctx.Done():
				log.Printf("Stopping kommence.\n", color.Bold)
				r.Stop(ctx)
				os.Exit(0)
			}
		}

	},
}

// Completer for autocomplete in interactive mode.
type Completer = func(in prompt.Document) []prompt.Suggest

func executableCompleter(c *configuration.Configuration) Completer {
	var s []prompt.Suggest
	for _, e := range c.Execs.Commands {
		s = append(s, prompt.Suggest{
			Text:        e.ID,
			Description: e.Description,
		})
		if e.Shortcut != "" {
			s = append(s, prompt.Suggest{
				Text:        e.Shortcut,
				Description: e.Description,
			})
		}
	}
	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
}

func getExecutables(ctx context.Context, c *configuration.Configuration) string {
	return prompt.Input(">>> ", executableCompleter(c),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(b *prompt.Buffer) {
				ctx.Done()
				os.Exit(0) // log.Fatal doesn't work, but panic somehow avoids this issue...
			}}))
}

func podCompleter(c *configuration.Configuration) Completer {
	var s []prompt.Suggest
	for _, e := range c.Pods.Pods {
		s = append(s, prompt.Suggest{
			Text:        e.ID,
			Description: e.Description,
		})
	}
	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
}

func getPods(ctx context.Context, c *configuration.Configuration) string {
	return prompt.Input(">>> ", podCompleter(c),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(b *prompt.Buffer) {
				ctx.Done()
				os.Exit(0) // log.Fatal doesn't work, but panic somehow avoids this issue...
			}}))
}

func flowCompleter(c *configuration.Configuration) Completer {
	var s []prompt.Suggest
	for _, e := range c.Flows.Flows {
		s = append(s, prompt.Suggest{
			Text:        e.ID,
			Description: e.Description,
		})
	}
	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
}

func getFlows(ctx context.Context, c *configuration.Configuration) string {
	return prompt.Input(">>> ", flowCompleter(c),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(b *prompt.Buffer) {
				ctx.Done()
				os.Exit(0)
			}}))
}

func startInteractiveExecutables(ctx context.Context, log *output.Logger, c *configuration.Configuration) []string {
	log.Printf("Select executables to run then press Enter:\n", color.Bold)
	var execs []string
	valid := false
	msg := ""
	for !valid {
		log.Printf("Available: %v\n", strings.Join(c.ListExecutables(), ", "), color.Bold)
		in := getExecutables(ctx, c)
		if in == "" {
			valid = true
		} else {
			execs = strings.Split(in, " ")
			valid, msg = c.ValidExecutables(execs)
			if !valid {
				log.Printf(msg+"\n", color.Bold)
			}
		}
	}
	return execs
}

func startInteractivePods(ctx context.Context, log *output.Logger, c *configuration.Configuration) []string {
	log.Printf("Select pods to forward then press Enter:\n", color.Bold)
	var pods []string
	valid := false
	msg := ""
	for !valid {
		log.Printf("Available: %v\n", strings.Join(c.ListPods(), ", "), color.Bold)
		in := getPods(ctx, c)
		if in == "" {
			valid = true
		} else {
			pods = strings.Split(in, " ")
			valid, msg = c.ValidPods(pods)
			if !valid {
				log.Printf(msg+"\n", color.Bold)
			}
		}
	}
	return pods
}

func startInteractiveFlow(ctx context.Context, log *output.Logger, c *configuration.Configuration) ([]string, []string) {
	log.Printf("Select flows to run then press Enter:\n", color.Bold)
	var execs []string
	var pods []string

	ex := make(map[string]bool)
	po := make(map[string]bool)

	valid := false
	msg := ""
	for !valid {
		log.Printf("Available: %v\n", strings.Join(c.ListFlows(), ", "), color.Bold)
		in := getFlows(ctx, c)
		if in == "" {
			valid = true
		} else {
			flows = strings.Split(in, " ")
			valid, msg = c.ValidFlows(flows)
			if !valid {
				log.Printf(msg+"\n", color.Bold)
			}
			for _, flow := range flows {
				theseExecs := c.Flows.GetExecutables(flow)
				for _, e := range theseExecs {
					ex[e] = true
				}
				thesePods := c.Flows.GetPods(flow)
				for _, p := range thesePods {
					po[p] = true
				}
			}
		}
	}
	for e := range ex {
		execs = append(execs, e)
	}
	for p := range po {
		pods = append(pods, p)
	}
	return execs, pods
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func startInteractive(ctx context.Context, log *output.Logger, c *configuration.Configuration) (*runner.Runner, *runner.Runtime) {
	var execs []string
	var pods []string
	if interactiveExecs && len(c.ListExecutables()) > 0 {
		execs = startInteractiveExecutables(ctx, log, c)
	}
	if interactivePods && len(c.ListPods()) > 0 {
		pods = startInteractivePods(ctx, log, c)
	}

	if interactiveFlows && len(c.ListFlows()) > 0 {
		otherExecs, otherPods := startInteractiveFlow(ctx, log, c)
		for _, exec := range otherExecs {
			if !contains(execs, exec) {
				execs = append(execs, exec)
			}
		}
		for _, pod := range otherPods {
			if !contains(pods, pod) {
				pods = append(pods, pod)
			}
		}
	}

	r := runner.New(log, c)
	return r, &runner.Runtime{Executables: execs, Pods: pods}

}

func startCommandLine(ctx context.Context, log *output.Logger, c *configuration.Configuration) (*runner.Runner, *runner.Runtime) {
	r := runner.New(log, c)
	for _, flow := range flows {
		otherExecs, otherPods := c.Flows.GetExecutables(flow), c.Flows.GetPods(flow)
		for _, exec := range otherExecs {
			if !contains(execs, exec) {
				execs = append(execs, exec)
			}
		}
		for _, pod := range otherPods {
			if !contains(pods, pod) {
				pods = append(pods, pod)
			}
		}
	}
	return r, &runner.Runtime{Executables: execs, Pods: pods}
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Executables
	startCmd.PersistentFlags().BoolVarP(&interactiveExecs, "interactive_execs", "X", false, "Interactive mode for executables")
	startCmd.PersistentFlags().StringSliceVarP(&execs, "execs", "x", nil, "Executables to run")

	// Pods
	startCmd.PersistentFlags().BoolVarP(&interactivePods, "interactive_pods", "P", false, "Interactive mode for pods")
	startCmd.PersistentFlags().StringSliceVarP(&pods, "pods", "p", nil, "Pods to forward")

	// Flows
	startCmd.PersistentFlags().BoolVarP(&interactiveFlows, "interactive_flows", "F", false, "Interactive mode for flows")
	startCmd.PersistentFlags().StringSliceVarP(&flows, "flows", "f", nil, "Pods to forward")
}
