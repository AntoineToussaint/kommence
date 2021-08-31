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

var interactive bool

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
		log := output.NewLogger(debug)
		config, err := configuration.Load(log, kommenceDir)
		if err != nil {
			log.Errorf(err.Error(), color.FgRed, color.Bold)
			os.Exit(1)
		}

		if interactive {
			startInteractive(log, config)
		}
	},
}

type Completer = func(in prompt.Document) []prompt.Suggest

func executableCompleter(c *configuration.Configuration) Completer {
	var s []prompt.Suggest
	for _, e := range c.Execs.Commands {
		s = append(s, prompt.Suggest{
			Text:        e.Name,
			Description: e.Description,
		})
	}
	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
}

func getRuns(ctx context.Context, c *configuration.Configuration) string {
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

func startInteractive(log *output.Logger, c *configuration.Configuration) {
	ctx := context.Background()
	log.Printf("Select executables to run then press Enter:\n")

	in := getRuns(ctx, c)
	runs := strings.Split(in, " ")
	r := runner.New(log, c)

	cancel := make(chan os.Signal, 2)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := r.Run(ctx, runner.Configuration{Runs: runs})
		if err != nil {
			log.Errorf("can't run")
		}
	}()
	<-cancel
	ctx.Done()
	os.Exit(1)

}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
}
