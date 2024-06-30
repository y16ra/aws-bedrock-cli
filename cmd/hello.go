/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/y16ra/aws-bedrock-cli/bedrock"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "say hello to Anthropic Claude",
	Long:  `This command sends a prompt to Anthropic Claude and prints the response.`,
	Run: func(cmd *cobra.Command, args []string) {
		bedrock.Hello(cmd.Context(), prompt, region)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVarP(&prompt, "prompt", "p", "Hello, how are you today?", "The prompt to send to Anthropic Claude")
	helloCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "The region to send the request to")
}
