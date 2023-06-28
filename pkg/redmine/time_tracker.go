package redmine

import (
	"context"
	"encoding/json"
	"net/http"
)

type TimeEntryRequest struct {
	issueID int       `json:"issue_id"`
	entry   TimeEntry `json:"time_entry"`
}
type TimeEntry struct {
	IssueID    int    `json:"issue_id"`
	SpentOn    string `json:"spent_on"`
	Hours      string `json:"hours"`
	ActivityID string `json:"activity_id"`
	Comments   string `json:"comments"`
	//	UserId      string            `json:"user_id"`
	CustomField CustomFieldObject `json:"custom_fields"`
}
type CustomFieldObject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CustomFieldPossibleValueObject struct {
	Value string `json:"value"`
}
func (c *Redmine) TimeTracker(ctx context.Context, SpendTime TimeEntryRequest) (err error) {
	path := "/time_entries.json"
	var body []byte
	body, err = json.Marshal(SpendTime)
	if err != nil {
		println("Cannot Marshal SpendTime")
	}

	err = c.Request(ctx, http.MethodPost, path, body, nil, nil)
	return
}
