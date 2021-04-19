package cmd

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runs []string

// flowCmd represents the flow command
var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config, _ := pkg.LoadFlowConfiguration(viper.GetViper())
		pkg.Flow(ctx, config, pkg.FlowRuntime{Runs: runs})
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.PersistentFlags().StringSliceVar(&runs, "run", nil, "environment")

}
