package mattermost

import (
	"errors"

	"mattermost/pkg/bun"
	"mattermost/pkg/redmine"

	mt "github.com/mattermost/mattermost-server/v6/model"
)

type Client struct {
	// mattermost client
	client *mt.Client4

	// webSocket client
	ws *mt.WebSocketClient

	// bot user instance
	botuser *mt.User

	// mattermost internal configurations.
	config *Config

	db *bun.DB

	redmine *redmine.Redmine
}

// New connects and init mattermost client.
func New(conf *Config,db *bun.DB) (c *Client, err error) {

	c = new(Client)

	// connect to mattermost client.
	c.client = mt.NewAPIv4Client(conf.ServerAddress)

	// login bot used by mattermost client.
	var resp *mt.Response
	if c.botuser, resp, err = c.client.Login(conf.UserEmail, conf.UserPassword); err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("error connecting to mattermost client")
	}

	// update bot username
	c.botuser.Username = conf.UserName

	// connect to websocket used by mattermost client.
	if c.ws, err = mt.NewWebSocketClient4(conf.ServerWsAddress, c.client.AuthToken); err != nil {
		return nil, err
	}

	c.config = conf

	c.db = db
	return
}

// find team by bot name.
func (c *Client) GetTeamByName() (err error) {

	if _, _, err = c.client.GetTeamByName(c.config.TeamName, ""); err != nil {
		return
	}

	return
}
