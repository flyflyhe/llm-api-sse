package api

import (
	"bm/internal/form"
	"bm/pkg/ai"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sse"
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

	ds := ai.NewDs("76227cb5-a62e-4f1b-83ec-20f67137442f")

	rc.SetStatusCode(http.StatusOK)
	s := sse.NewStream(rc)

	stream, err := ds.ChatCompletions(reqForm)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)

		h.Fail(rc, err.Error())
		return
	}

	defer stream.Close()

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
			fmt.Print(recv.Choices[0].Delta.Content)
		}
	}

	if err = s.Publish(&sse.Event{Event: "done"}); err != nil {
		return
	}
}
