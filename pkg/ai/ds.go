package ai

import (
	ark "github.com/sashabaranov/go-openai"
	"sync"
)

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

func (d *Ds) ChatCompletions() {
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
