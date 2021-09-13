package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/antoinetoussaint/kommence/pkg/runner"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var interactive bool
var execs []string
var pods []string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		if interactive {
			log.Debugf("starting interactive mode\n")
			startInteractive(ctx, log, config)
		} else if len(execs) > 0 || len(pods) > 0 {
			log.Debugf("starting runner mode\n")
			startRunner(ctx, log, config)
		}

		<-cancel
		stop()
		log.Printf("\nStopping all running processes.\n", color.Bold)
		time.Sleep(time.Second)
		os.Exit(0)
	},
}

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

func startInteractive(ctx context.Context, log *output.Logger, c *configuration.Configuration) {
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
				log.Printf(msg + "\n", color.Bold)
			}
		}
	}

	log.Printf("Select pods to forward then press Enter:\n", color.Bold)
	
	var pods []string
	valid = false
	msg = ""
	for !valid {
		log.Printf("Available: %v\n", strings.Join(c.ListPods(), ", "), color.Bold)
		in := getPods(ctx, c)
		if in == "" {
			valid = true
		} else {
			pods = strings.Split(in, " ")
			valid, msg = c.ValidPods(execs)
			if !valid {
				log.Printf(msg + "\n", color.Bold)
			}
		}
	}

	r := runner.New(log, c)

	go func() {
		err := r.Run(ctx, runner.Configuration{Executables: execs, Pods: pods})
		if err != nil {
			log.Errorf("can't run")
		}
	}()

}

func startRunner(ctx context.Context, log *output.Logger, c *configuration.Configuration) {
	r := runner.New(log, c)

	go func() {
		err := r.Run(ctx, runner.Configuration{Executables: execs, Pods: pods})
		if err != nil {
			log.Errorf("can't run")
		}
	}()

}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
	startCmd.PersistentFlags().StringSliceVarP(&execs, "execs", "x", nil, "Executables to run")
	startCmd.PersistentFlags().StringSliceVarP(&pods, "pods", "p", nil, "Pods to forward")
}
