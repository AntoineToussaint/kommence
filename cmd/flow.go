package cmd

import (
	"context"
	"fmt"
	"github.com/antoinetoussaint/kommence/pkg"
	"github.com/spf13/cobra"
)

var runs []string
var kubes []string

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
		config, err := pkg.LoadFlowConfiguration(kommenceDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		pkg.Flow(ctx, config, pkg.FlowRuntime{Runs: runs, Kubes: kubes})
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.PersistentFlags().StringSliceVar(&runs, "run", nil, "environment")
	flowCmd.PersistentFlags().StringSliceVar(&kubes, "kube", nil, "environment")
}
