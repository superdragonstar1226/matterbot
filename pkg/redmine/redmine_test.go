package redmine

import (
	"context"
	"net/http"
	"testing"

	"mattermost/pkg/logger"

	"github.com/stretchr/testify/require"
)

// TODO: implement tests for redmine methods.

func NewTestClient(t *testing.T) *Redmine {
	config := new(Config)
	config.Host = "https://redmine.sfxdx.com"
	config.APIKey = "e4a6852ceeb126a9a7d387d394c44dd4210da02b"
	redmine := New(logger.NewNop(), config)
	return redmine
}
func TestRedmine_Init(t *testing.T) {
	redmine := NewTestClient(t)
	require.NotEmpty(t, redmine)
	require.Equal(t, "https://redmine.sfxdx.com", redmine.config.Host)
	require.Equal(t, "e4a6852ceeb126a9a7d387d394c44dd4210da02b", redmine.config.APIKey)
	require.Equal(t, logger.NewNop(), redmine.log)
}

func TestRedmine_RequestURI(t *testing.T) {
	redmine := NewTestClient(t)
	var (
		uri string
		err error
	)

	type query struct {
		Val1 string `url:"val1"`
		Val2 string `url:"val2"`
	}
	var (
		testPath  = "testpath"
		testQuery = &query{
			Val1: "param1",
			Val2: "param2",
		}
	)
	uri, err = redmine.RequestURI(testPath, testQuery)
	require.NoError(t, err)
	require.Equal(t,
		"https://redmine.sfxdx.comtestpath?val1=param1&val2=param2",
		uri)
}

func TestRedmine_Request(t *testing.T) {
	ctx := context.Background()
	c := NewTestClient(t)
	var l *ListResponse

	err := c.Request(ctx, http.MethodGet, "issues.json", nil, nil, l)
	require.Equal(t, nil, err)
	require.Equal(t, 1, l.TotalCount)
	require.Equal(t, 9984, l.Issues[0].ID)

}
