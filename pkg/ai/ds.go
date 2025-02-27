package ai

import (
	"bm/internal/form"
	"context"
	ark "github.com/sashabaranov/go-openai"
	"sync"
)

var msgStore = map[string][]ark.ChatCompletionMessage{}
var msgMu sync.Mutex

type Ds struct {
	token string

	mu     sync.Mutex
	client *ark.Client
}

func NewDs(token string) *Ds {
	return &Ds{
		token: token,
		mu:    sync.Mutex{},
	}
}

func (d *Ds) FindOne(chatSessionId string) []ark.ChatCompletionMessage {
	return msgStore[chatSessionId]
}

func (d *Ds) SetMessages(chatSessionId string, messages []ark.ChatCompletionMessage) {
	msgMu.Lock()
	defer msgMu.Unlock()
	msgStore[chatSessionId] = messages
}

func (d *Ds) ChatCompletions(reqForm form.DsRequest) (stream *ark.ChatCompletionStream, err error) {
	client := d.GetClient()

	msgList := d.FindOne(reqForm.ChatSessionId)
	if msgList == nil {
		msgList = []ark.ChatCompletionMessage{
			{
				Role:    ark.ChatMessageRoleSystem,
				Content: reqForm.Prompt,
			},
		}
	} else {
		msgList = append(msgList, ark.ChatCompletionMessage{
			Role:    ark.ChatMessageRoleUser,
			Content: reqForm.Prompt,
		})
	}

	d.SetMessages(reqForm.ChatSessionId, msgList)

	stream, err = client.CreateChatCompletionStream(
		context.Background(),
		ark.ChatCompletionRequest{
			Model:    "ep-20250220171927-7w4pc",
			Messages: msgList,
		},
	)

	return
}

func (d *Ds) GetClient() *ark.Client {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.client != nil {
		return d.client
	}

	config := ark.DefaultConfig("76227cb5-a62e-4f1b-83ec-20f67137442f")
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	d.client = ark.NewClientWithConfig(config)

	return d.client
}
