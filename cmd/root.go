package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	prompt string
	region string
)

var rootCmd = &cobra.Command{
	Use:   "aws-bedrock-cli",
	Short: "A CLI tool for Anthropic Claude",
	Long:  `This CLI tool is a collection of commands to interact with Anthropic Claude.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Hugo is a very fast static site generator")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
