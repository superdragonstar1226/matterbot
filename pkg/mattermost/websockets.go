package mattermost

import (
	// "encoding/json"
	"context"

	"github.com/mattermost/mattermost-server/v6/model"
)

// subccribe

// handle wss connection

// dial wss

// close wss connection

// handle wss event

// Listen ws client.
// TODO: need to some refactor.

type ReportForRedmine struct {
	IssueNumber string
	Hours       string
	Activity    string
	ReportText  string
}

func (c *Client) DialWebSockets(ctx context.Context) {
	// do hard job
	c.ws.Listen()
	go func() {
		for resp := range c.ws.EventChannel {
			c.handleWebSocketResponse(ctx,resp)
		}
	}()

	select {}
}

func (c *Client) handleWebSocketResponse(ctx context.Context,event *model.WebSocketEvent) {
	if event.EventType() == model.WebsocketEventPosted {
		c.handlMsg(ctx, event)
	}
}
