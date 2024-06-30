package bedrock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/y16ra/aws-bedrock-cli/model"
)

type InvokeModelWithResponseStreamWrapper struct {
	BedrockRuntimeClient *bedrockruntime.Client
}

// claude3 request data type
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type RequestBodyClaude3 struct {
	MaxTokensToSample int       `json:"max_tokens"`
	Temperature       float64   `json:"temperature,omitempty"`
	AnthropicVersion  string    `json:"anthropic_version"`
	Messages          []Message `json:"messages"`
}

// claude3 response data type
type Delta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseClaude3 struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta Delta  `json:"delta"`
}

type Response struct {
	Completion string `json:"completion"`
}

func (wrapper InvokeModelWithResponseStreamWrapper) InvokeModelWithResponseStream(ctx context.Context, prompt string) (string, error) {

	modelID := "anthropic.claude-v2:1"

	// Anthropic Claude requires you to enclose the prompt as follows:
	prefix := "Human: "
	postfix := "\n\nAssistant:"
	prompt = prefix + prompt + postfix

	request := model.ClaudeRequest{
		Prompt:            prompt,
		MaxTokensToSample: 1000,
		Temperature:       0.5,
		StopSequences:     []string{"\n\nHuman:"},
	}

	body, err := json.Marshal(request)
	if err != nil {
		log.Panicln("Couldn't marshal the request: ", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModelWithResponseStream(context.Background(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		Body:        body,
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "no such host") {
			log.Printf("The Bedrock service is not available in the selected region. Please double-check the service availability for your region at https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/.\n")
		} else if strings.Contains(errMsg, "Could not resolve the foundation model") {
			log.Printf("Could not resolve the foundation model from model identifier: \"%v\". Please verify that the requested model exists and is accessible within the specified region.\n", modelID)
		} else {
			log.Printf("Couldn't invoke Anthropic Claude. Here's why: %v\n", err)
		}
	}
	resp, err := processStreamingOutput(output, func(ctx context.Context, part []byte) error {
		fmt.Print(string(part))
		return nil
	})

	if err != nil {
		log.Fatal("streaming output processing error: ", err)
	}

	return resp.Completion, nil

}

type StreamingOutputHandler func(ctx context.Context, part []byte) error

func processStreamingOutput(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (Response, error) {

	var combinedResult string
	resp := Response{}

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:

			//fmt.Println("payload", string(v.Value.Bytes))

			var resp Response
			// var resp ResponseClaude3
			err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&resp)
			if err != nil {
				return resp, err
			}

			err = handler(context.Background(), []byte(resp.Completion))
			if err != nil {
				return resp, err
			}

			combinedResult += resp.Completion

		case *types.UnknownUnionMember:
			fmt.Println("unknown tag:", v.Tag)

		default:
			fmt.Println("union is nil or unknown type")
		}
	}

	resp.Completion = combinedResult

	return resp, nil
}
