package cmd

import (
	"github.com/AntoineToussaint/jarvis/pkg/configuration"
	"github.com/AntoineToussaint/jarvis/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all possible tasks to run",
	Run: func(cmd *cobra.Command, args []string) {
		log := output.NewLogger(debug)
		config, err := configuration.Load(log, jarvisDir)
		if err != nil {
			log.Errorf(err.Error()+"\n", color.FgRed, color.Bold)
			os.Exit(1)
		}
		config.Print(log)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
