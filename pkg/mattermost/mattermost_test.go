package mattermost

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func connect(t *testing.T) (*Client, error) {

// 	conf := &Config{
// 		UserPassword:           "qwerty",
// 		UserEmail:              "aatrushkevich+3@sfxdx.com",
// 		UserName:               "Testuser4",
// 		TeamName:               "Test",
// 		ServerAddress:          "http://localhost:8065",
// 		ServerWsAddress:        "WS://localhost:8065",
// 		TimeBeforeNotification: 0,
// 		ReportFileName:         "",
// 	}

// 	cli, err := New(conf)
// 	require.NoError(t, err)

// 	return cli, err
// }

// func TestClient_Connect(t *testing.T) {
// 	client, err := connect(t)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, client)
// }
//