/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/spf13/cobra"
)

const (
	KNOWLEDGE_BASE_MODEL_ID         = "anthropic.claude-3-haiku-20240307-v1:0"
	KNOWLEDGE_BASE_NUMBER_OF_RESULT = 6
)

var knowledgeBaseID string

// knowledgeCmd represents the knowledge command
var knowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "Generates answers using Anthropic Claude with knowledge base information",
	Long: `This command generates answers using Anthropic Claude with knowledge base information.
		You need to provide the knowledge base ID to retrieve the knowledge base.`,
	Run: func(cmd *cobra.Command, args []string) {
		// invoke bedrock agent runtime to retrieve opensearch
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
		if err != nil {
			fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
			fmt.Println(err)
			return
		}
		client := bedrockagentruntime.NewFromConfig(cfg)
		output, err := client.RetrieveAndGenerate(
			context.TODO(),
			&bedrockagentruntime.RetrieveAndGenerateInput{
				Input: &types.RetrieveAndGenerateInput{
					Text: aws.String(prompt),
				},
				RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
					Type: types.RetrieveAndGenerateTypeKnowledgeBase,
					KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
						KnowledgeBaseId: aws.String(knowledgeBaseID),
						ModelArn:        aws.String(KNOWLEDGE_BASE_MODEL_ID),
						RetrievalConfiguration: &types.KnowledgeBaseRetrievalConfiguration{ // optional
							VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
								NumberOfResults: aws.Int32(KNOWLEDGE_BASE_NUMBER_OF_RESULT),
							},
						},
					},
				},
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Results: %v\n", *output.Output.Text)
	},
}

func init() {
	rootCmd.AddCommand(knowledgeCmd)
	knowledgeCmd.Flags().StringVarP(&prompt, "prompt", "p", "Hello, how are you today?", "The prompt to send to Anthropic Claude")
	knowledgeCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "The region to send the request to")
	knowledgeCmd.Flags().StringVarP(&knowledgeBaseID, "knowledgeBaseID", "k", "", "The knowledge base ID")
	knowledgeCmd.MarkFlagRequired("knowledgeBaseID")
}
