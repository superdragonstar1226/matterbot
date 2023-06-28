package mattermost

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

func (c *Client) sendMsgToUser(msg string, ChannelId string) {
	post := &model.Post{
		ChannelId: ChannelId,
		Message:   msg,
		UserId:    c.botuser.Id,
	}
	if _, _, err := c.client.CreatePost(post); err != nil {
		println("We failed to send a message to the user")

	}
}
