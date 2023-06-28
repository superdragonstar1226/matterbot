package redmine

import (
	"context"
	"testing"

	// "mattermost/pkg/logger"

	"github.com/stretchr/testify/require"
)

func TestRedmine_Retrieve(t *testing.T) {
	ctx := context.Background()
	r := &RetrieveRequest{
		IssueID: 9984,
	}
	c := NewTestClient(t)
	result, err := c.Retrieve(ctx, r)
	require.Equal(t, nil, err)
	require.NotEmpty(t, result)
	require.Equal(t, 687, result.ProjectID)
	require.Equal(t, "test-feature-1", result.Subject)

}

// func TestRedMine_ToStopped(t *testing.T) {
// 	ctx := context.Background()
// 	c := NewTestClient(t)
// 	r := &ToStoppedRequest{
// 		UsedID: 266,
// 	}
// 	result := c.ToStopped(ctx, r)
// 	require.Equal(t,

// }

// func TestRedmine_ToProgress(t *testing.T) {
// 	ctx := context.Background()
// 	c := NewTestClient(t)

// }

// func TestRedmine_List(t *testing.T) {
// 	ctx := context.Background()
// 	c := NewTestClient(t)
// }
