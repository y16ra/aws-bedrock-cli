/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/spf13/cobra"
	"github.com/y16ra/aws-bedrock-cli/bedrock"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(cmd.Context(), config.WithRegion(region))
		if err != nil {
			log.Fatalf("unable to load aws config, %v", err)
		}
		client := bedrockruntime.NewFromConfig(cfg)
		responseStreamWrapper := bedrock.InvokeModelWithResponseStreamWrapper{
			BedrockRuntimeClient: client,
		}
		responseStreamWrapper.InvokeModelWithResponseStream(cmd.Context(), prompt)
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringVarP(&prompt, "prompt", "p", "Hello, how are you today?", "The prompt to send to Anthropic Claude")
	streamCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "The region to send the request to")
}
