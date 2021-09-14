package cmd

import (
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
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
			log.Errorf(err.Error()+"\n", color.FgRed, color.Bold)
			os.Exit(1)
		}
		config.Print(log)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
