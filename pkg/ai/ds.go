package ai

import (
	"bm/internal/config"
	"bm/internal/form"
	"bm/internal/tool"
	"bm/pkg/logging"
	"context"
	"encoding/json"
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

func (d *Ds) TestFunction() ([]ark.ChatCompletionMessage, error) {
	param := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]string{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
		},
		"required": []string{"location"},
	}

	pj, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	logging.Logger.WithCtx(d.ctx).Info(string(pj))
	tools := []ark.Tool{{
		Type: "function",
		Function: &ark.FunctionDefinition{
			Name:        "get_weather",
			Description: "Get weather of an location, the user shoud supply a location first",
			Strict:      false,
			Parameters:  param,
		},
	}}

	var messageList []ark.ChatCompletionMessage
	message := ark.ChatCompletionMessage{
		Role:    "user",
		Content: "How's the weather in Hangzhou?",
	}
	messageList = append(messageList, message)

	client := d.GetClient()

	res, err := client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model:    d.model,
			Messages: messageList,
			Tools:    tools,
		},
	)

	if err != nil {
		logging.Logger.WithCtx(d.ctx).Error("ChatCompletions", err)
		return nil, err
	}

	logging.Logger.WithCtx(d.ctx).Info("ChatCompletions", res.Choices[0].Message)

	messageRes := res.Choices[0].Message
	toolCall := messageRes.ToolCalls[0]

	//调用自己的真实函数

	messageList = append(messageList, messageRes)
	messageList = append(messageList, ark.ChatCompletionMessage{
		Role:         ark.ChatMessageRoleTool,
		Content:      "24℃",
		Refusal:      "",
		MultiContent: nil,
		Name:         "",
		FunctionCall: nil,
		ToolCalls:    nil,
		ToolCallID:   toolCall.ID,
	})

	res, err = client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model:    d.model,
			Messages: messageList,
			Tools:    tools,
		},
	)

	if err != nil {
		logging.Logger.WithCtx(d.ctx).Error("ChatCompletions", err)
		return nil, err
	}

	messageList = append(messageList, res.Choices[0].Message)
	resJ, _ := json.Marshal(res)
	logging.Logger.WithCtx(d.ctx).Info("ChatCompletions", string(resJ))
	return messageList, nil
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
