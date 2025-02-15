package ai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Client struct {
	*openai.Client
}

func NewClient(endpoint, token string) *Client {
	client := openai.NewClient(option.WithBaseURL(endpoint), option.WithAPIKey(token))
	return &Client{
		client,
	}

}

func (c *Client) ChatCompletion(ctx context.Context, req string) (resp string, err error) {
	completion, err := c.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(req),
		}),
		Seed: openai.Int(1),
		// Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return "", err
	}
	return completion.Choices[0].Message.Content, nil
}

func Hello() {
	opts := option.WithBaseURL("http://127.0.0.1:8000/")

	client := openai.NewClient(opts)

	ctx := context.Background()

	question := "напиши хокку"

	print("> ")
	println(question)
	println()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed: openai.Int(1),
		// Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err)
	}

	// fmt.Println(completion)

	println(completion.Choices[0].Message.Content)
}
