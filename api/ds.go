package api

import (
	"bm/pkg/ai"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	ark "github.com/sashabaranov/go-openai"
	"io"
)

type Ds struct {
	Base
}

func (h *Ds) ChatCompletion(ctx context.Context, rc *app.RequestContext) {
	ds := ai.NewDs("76227cb5-a62e-4f1b-83ec-20f67137442f")

	rc.Response.Header.Set("Content-Type", "text/event-stream")
	rc.Response.Header.Set("Cache-Control", "no-cache")
	rc.Response.Header.Set("Connection", "keep-alive")

	client := ds.GetClient()

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		ark.ChatCompletionRequest{
			Model: "ep-20250220171927-7w4pc",
			Messages: []ark.ChatCompletionMessage{
				{
					Role:    ark.ChatMessageRoleSystem,
					Content: "你是人工智能助手",
				},
				{
					Role:    ark.ChatMessageRoleUser,
					Content: "你能做什么?",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)

		h.Fail(rc, err.Error())
		return
	}

	defer stream.Close()

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Printf("Stream chat error: %v\n", err)
			return
		}

		if len(recv.Choices) > 0 {
			fmt.Print(recv.Choices[0].Delta.Content)
		}
	}
}
