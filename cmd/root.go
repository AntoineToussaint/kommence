package cmd

import (
	"fmt"
	"github.com/AntoineToussaint/jarvis/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var jarvisDir string
var kubeConfigPath string
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jarvis",
	Short: "Start multiple tasks: executables, pod forwarding, flows",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log := output.NewLogger(debug)
		log.Printf("Welcome to jarvis!\n", color.Bold)
		log.Printf("To get started, run: ")
		log.Printf(" jarvis init\n", color.Bold)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&jarvisDir, "config", "jarvis", "jarvis folder (default is jarvis")
	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kube", "", "kubernetes config path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
