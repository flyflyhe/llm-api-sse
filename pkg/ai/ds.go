package ai

import (
	"bm/internal/config"
	"bm/internal/form"
	"bm/internal/tool"
	"bm/pkg/logging"
	"context"
	ark "github.com/sashabaranov/go-openai"
	"sync"
)

var msgStore = map[string][]ark.ChatCompletionMessage{}
var msgMu sync.Mutex

type Ds struct {
	token   string
	baseUrl string
	model   string

	mu     sync.Mutex
	client *ark.Client
	ctx    context.Context
}

func NewDs(ctx context.Context, llmConfig config.LLM) *Ds {
	return &Ds{
		token:   llmConfig.Token,
		baseUrl: llmConfig.BaseUrl,
		model:   llmConfig.Model,
		ctx:     ctx,
		mu:      sync.Mutex{},
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

func (d *Ds) AppendMessageV1(chatSessionId string, message string) {
	msgMu.Lock()
	defer msgMu.Unlock()

	role := ark.ChatMessageRoleUser
	if msgStore[chatSessionId] == nil {
		role = ark.ChatMessageRoleSystem
	}

	msgStore[chatSessionId] = append(msgStore[chatSessionId], ark.ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
}

func (d *Ds) AppendMessage(chatSessionId string, message, role string) {
	msgMu.Lock()
	defer msgMu.Unlock()
	msgStore[chatSessionId] = append(msgStore[chatSessionId], ark.ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
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

	logging.Logger.WithCtx(d.ctx).Info("msgList", tool.ToJson(msgList))
	stream, err = client.CreateChatCompletionStream(
		context.Background(),
		ark.ChatCompletionRequest{
			Model:    d.model,
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

	arkConfig := ark.DefaultConfig(d.token)
	arkConfig.BaseURL = d.baseUrl
	d.client = ark.NewClientWithConfig(arkConfig)

	return d.client
}
