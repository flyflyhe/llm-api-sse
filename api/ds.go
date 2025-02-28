package api

import (
	"bm/internal/config"
	"bm/internal/form"
	"bm/pkg/ai"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sse"
	ark "github.com/sashabaranov/go-openai"
	"io"
	"net/http"
)

type Ds struct {
	Base
}

func (h *Ds) ChatCompletion(ctx context.Context, rc *app.RequestContext) {
	var reqForm form.DsRequest
	if err := rc.BindAndValidate(&reqForm); err != nil {
		h.Fail(rc, err.Error())
		return
	}

	ds := ai.NewDs(ctx, config.GetApp().LLM)

	rc.SetStatusCode(http.StatusOK)
	s := sse.NewStream(rc)

	stream, err := ds.ChatCompletions(reqForm)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)

		h.Fail(rc, err.Error())
		return
	}

	defer stream.Close()

	fullContent := ""
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Stream chat error: %v\n", err)
			return
		}

		if len(recv.Choices) > 0 {
			event := &sse.Event{
				Event: "message",
				Data:  []byte(recv.Choices[0].Delta.Content),
			}
			err := s.Publish(event)
			if err != nil {
				return
			}
			fullContent += recv.Choices[0].Delta.Content
		}
	}

	ds.AppendMessageV1(reqForm.ChatSessionId, reqForm.Prompt)
	ds.AppendMessage(reqForm.ChatSessionId, fullContent, ark.ChatMessageRoleAssistant)

	if err = s.Publish(&sse.Event{Event: "done"}); err != nil {
		return
	}
}

func (h *Ds) TestFunction(ctx context.Context, rc *app.RequestContext) {

	ds := ai.NewDs(ctx, config.GetApp().LLM)

	if res, err := ds.TestFunction(); err != nil {
		h.Fail(rc, err.Error())
		return
	} else {
		h.Success(rc, res)
	}
}
