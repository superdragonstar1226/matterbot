package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

const limitDefault = 100

// IssuesRedmine
type RedmineIssues struct {
	*Redmine
}

// IssuesInterface wrpas Issues methods.
type IssuesInterface interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
	// TODO: implement other methods here.
}

// IssuesInterface manager.
func (r *Redmine) Issues() *RedmineIssues {
	return &RedmineIssues{Redmine: r}
}

type IssueAllGetRequest struct {
	Includes []string
	Filters  IssueGetRequestFilters
}
type IssueGetRequestFilters struct {
	Fields map[string][]string
	Cf     []IssueGetRequestFiltersCf
}
type IssueGetRequestFiltersCf struct {
	ID    int
	Value string
}
type IssueResult struct {
	Issues     []IssueObject `json:"issues"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

type IDName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IssueObject struct {
	ID             int                    `json:"id"`
	Project        IDName                 `json:"project"`
	Tracker        IDName                 `json:"tracker"`
	Status         IDName                 `json:"status"`
	Priority       IDName                 `json:"priority"`
	Author         IDName                 `json:"author"`
	AssignedTo     IDName                 `json:"assigned_to"`
	Category       IDName                 `json:"category"`
	FixedVersion   IDName                 `json:"fixed_version"`
	Parent         IssueParentObject      `json:"parent"`
	Subject        string                 `json:"subject"`
	Description    string                 `json:"description"`
	StartDate      string                 `json:"start_date"`
	DueDate        string                 `json:"due_date"`
	DoneRatio      int                    `json:"done_ratio"`
	IsPrivate      bool                   `json:"is_private"`
	EstimatedHours float64                `json:"estimated_hours"`
	SpentHours     float64                `json:"spent_hours"` // used only: get single issue
	CustomFields   []CustomFieldGetObject `json:"custom_fields"`
	CreatedOn      string                 `json:"created_on"`
	UpdatedOn      string                 `json:"updated_on"`
	ClosedOn       string                 `json:"closed_on"`
	Children       []IssueChildrenObject  `json:"children"`
	Attachments    []AttachmentObject     `json:"attachments"` // used only: get single issue
	Relations      []IssueRelationObject  `json:"relations"`
	Changesets     []IssueChangesetObject `json:"changesets"` // used only: get single issue
	Journals       []IssueJournalObject   `json:"journals"`   // used only: get single issue
	Watchers       []IDName               `json:"watchers"`   // used only: get single issue
}
type IssueParentObject struct {
	ID int `json:"id"`
}
type CustomFieldGetObject struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Multiple bool   `json:"multiple"`
	Value    string `json:"value"`
}
type IssueChildrenObject struct {
	ID       int                   `json:"id"`
	Tracker  IDName                `json:"tracker"`
	Subject  string                `json:"subject"`
	Children []IssueChildrenObject `json:"children"`
}
type AttachmentObject struct {
	ID          int    `json:"id"`
	FileName    string `json:"filename"`
	FileSize    string `json:"filesize"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
	ContentURL  string `json:"content_url"`
	Author      IDName `json:"author"`
	CreatedOn   string `json:"created_on"`
}
type IssueRelationObject struct {
	ID           int    `json:"id"`
	IssueID      int    `json:"issue_id"`
	IssueToID    int    `json:"issue_to_id"`
	RelationType string `json:"relation_type"`
	Delay        int    `json:"delay"`
}
type IssueChangesetObject struct {
	Revision    string `json:"revision"`
	User        IDName `json:"user"`
	Comments    string `json:"comments"`
	CommittedOn string `json:"committed_on"`
}
type IssueJournalObject struct {
	ID        int                        `json:"id"`
	User      IDName                     `json:"user"`
	Notes     string                     `json:"notes"`
	CreatedOn string                     `json:"created_on"`
	Details   []IssueJournalDetailObject `json:"details"`
}
type IssueJournalDetailObject struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

type IssueSingleGetRequest struct {
	Includes []string
}
type issueSingleResult struct {
	Issue IssueObject `json:"issue"`
}
type IssueMultiGetRequest struct {
	Includes []string
	Filters  IssueGetRequestFilters
	Offset   int
	Limit    int
}

type IssueUpdateObject struct {
	ProjectID      int                       `json:"project_id,omitempty"`
	TrackerID      int                       `json:"tracker_id,omitempty"`
	StatusID       int                       `json:"status_id,omitempty"`
	PriorityID     int                       `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject,omitempty"`
	Description    string                    `json:"description,omitempty"`
	CategoryID     int                       `json:"category_id,omitempty"`
	FixedVersionID int                       `json:"fixed_version_id,omitempty"`
	AssignedToID   int                       `json:"assigned_to_id,omitempty"`
	ParentIssueID  int                       `json:"parent_issue_id,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
	Notes          string                    `json:"notes,omitempty"`
}
type CustomFieldUpdateObject struct {
	ID    int         `json:"id"`
	Value interface{} `json:"value"` // can be a string or strings slice
}
type AttachmentUploadObject struct {
	ID          int    `json:"id,omitempty"`
	Token       string `json:"token"`
	Filename    string `json:"filename"`     // This field fills in AttachmentUpload() function, not by Redmine. User can redefine this value manually
	ContentType string `json:"content_type"` // This field fills in AttachmentUpload() function, not by Redmine. User can redefine this value manually
}

// ListRequest is a List method request.
type ListRequest struct {
	Includes []string
	Filters  IssueGetRequestFilters
	Offset   int
	Limit    int
}

// ListResponse is a List method response.
type ListResponse struct {
	Issues     []IssueObject `json:"issues"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

// List method returns the list of fetched Issues with optional filter params.
func (c *Redmine) List(ctx context.Context, r *ListRequest) (out *ListResponse, err error) {
	type query struct {
		Offset int `url:"offset"`
		Limit  int `url:"limit"`
		// TODO: add other url params
		//Includes []string `url:includes`
	}

	in := &query{
		Limit: r.Limit,
		// TODO: add other url params
	}

	out = new(ListResponse)
	err = c.Request(ctx, http.MethodGet, "/issues.json", nil, in, out)
	if err != nil {
		return
	}
	return
}

// RetrieveRequest is a Retrieve method request.
type RetrieveRequest struct {
	IssueID int
}

// RetrieveResponse is a Retrieve method request.
type RetrieveResponse struct {
	ProjectID      int                       `json:"project_id,omitempty"`
	TrackerID      int                       `json:"tracker_id,omitempty"`
	StatusID       int                       `json:"status_id,omitempty"`
	PriorityID     int                       `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject,omitempty"`
	Description    string                    `json:"description,omitempty"`
	CategoryID     int                       `json:"category_id,omitempty"`
	FixedVersionID int                       `json:"fixed_version_id,omitempty"`
	AssignedToID   int                       `json:"assigned_to_id,omitempty"`
	ParentIssueID  int                       `json:"parent_issue_id,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
	Notes          string                    `json:"notes,omitempty"`
}

//get single issue
func (c *Redmine) Retrieve(ctx context.Context, r *RetrieveRequest) (out *RetrieveResponse, err error) {
	path := "/issues/" + strconv.Itoa(r.IssueID) + ".json"

	out = new(RetrieveResponse)
	err = c.Request(ctx, http.MethodGet, path, nil, nil, out)
	if err != nil {
		return
	}

	return
}

type ToStoppedRequest struct {
	UsedID int
}

func (c *Redmine) ToStopped(ctx context.Context, r *ToStoppedRequest) (err error) {

	out := new(Issues)
	if err = c.Request(ctx, http.MethodGet, "/issues.json", nil, nil, out); err != nil {
		return
	}

	for _, v := range out.Issues {
		v.StatusID = 6
		var b []byte
		if b, err = json.Marshal(v); err != nil {
			return
		}

		if err = c.Request(ctx, http.MethodPut, "/issues/"+strconv.Itoa(r.UsedID)+".json", b, nil, nil); err != nil {
			return
		}
	}

	return
}

type ToProgressRequest struct {
	UserID  int
	IssueID int
}

type Issues struct {
	Issues []*IssueUpdateObject `json:"issues"`
}

type IssueUpdate struct {
	Issue IssueUpdateObject `json:"issue"`
}

func (c *Redmine) ToProgress(ctx context.Context, r *ToProgressRequest) (err error) {

	out := new(Issues)
	if err = c.Request(ctx, http.MethodGet, "/issues.json", nil, nil, out); err != nil {
		return
	}

	issue := new(IssueUpdate)
	issue.Issue.StatusID = 2
	var b []byte
	if b, err = json.Marshal(issue); err != nil {
		return
	}

	if err = c.Request(ctx, http.MethodPut, "/issues/"+strconv.Itoa(r.UserID)+".json", b, nil, nil); err != nil {
		return
	}

	return
}

// TODO: to delete deprecated method
// func (c *Redmine) IssuesAllGet(ctx context.Context, request IssueAllGetRequest) (IssueResult, error) {

// 	var (
// 		issues IssueResult
// 		offset int
// 	)

// 	m := IssueMultiGetRequest{
// 		Filters:  request.Filters,
// 		Includes: request.Includes,
// 		Limit:    limitDefault,
// 	}

// 	for {

// 		m.Offset = offset

// 		i, err := c.IssuesMultiGet(ctx, m)
// 		if err != nil {
// 			return issues, err
// 		}

// 		issues.Issues = append(issues.Issues, i.Issues...)

// 		if offset+i.Limit >= i.TotalCount {
// 			issues.TotalCount = i.TotalCount
// 			issues.Limit = i.TotalCount

// 			break
// 		}

// 		offset += i.Limit
// 	}
// 	return issues, nil
// }

// TODO: to delete deprecated method
// func (c *Redmine) IssuesMultiGet(ctx context.Context, request IssueMultiGetRequest) (IssueResult, error) {

// 	var i IssueResult

// 	urlParams := url.Values{

// 	}
// 	urlParams.Add("offset", strconv.Itoa(request.Offset))
// 	urlParams.Add("limit", strconv.Itoa(request.Limit))

// 	// Preparing includes
// 	urlIncludes(&urlParams, request.Includes)

// 	// Preparing filters
// 	issueURLFilters(&urlParams, request.Filters)

// 	path := "/issues.json"

// 	err := c.Get(ctx, &i, path)

// 	return i, err
// }

// TODO: to delete deprecated method
// func urlIncludes(urlParams *url.Values, includes []string) {

// 	if len(includes) == 0 {
// 		return
// 	}

// 	urlParams.Add("include", strings.Join(includes, ","))
// }

// TODO: to delete deprecated method
// func issueURLFilters(urlParams *url.Values, filters IssueGetRequestFilters) {

// 	// Filter fields (e.g. `issue_id`, `tracker_id`, etc)
// 	for n, s := range filters.Fields {
// 		urlParams.Add(n, strings.Join(s, ","))
// 	}

// 	// Custom fields
// 	for _, c := range filters.Cf {
// 		urlParams.Add("cf_"+strconv.Itoa(c.ID), c.Value)
// 	}
// }

// func (c *Redmine) IssueUpdate(ctx context.Context, id int, issue IssueUpdateObject) error {

// 	path := "/issues/" + strconv.Itoa(id) + ".json"

// 	err := c.Put(ctx, issueUpdate{Issue: issue}, nil, path)

// 	return err
// }

//проверить на нескольких задачах
// TODO: to depete deprecated method.
// func (c *Redmine) IssuesToStopped(ctx context.Context, userId int) error {
// 	// получить задачи которые в inprogress
// 	filter := IssueAllGetRequest{
// 		Filters: IssueGetRequestFilters{
// 			Fields: map[string][]string{
// 				"assigned_to_id": []string{strconv.Itoa(userId)},
// 				"status_id":      []string{"2"},
// 			},
// 		},
// 	}

// 	issues, err := c.IssuesAllGet(ctx, filter)
// 	if err != nil {
// 		return err
// 	}

// 	// если inprogress задач нет, то сразу на выход
// 	if issues.TotalCount == 0 { // нет задач в инпрогресс
// 		return nil
// 	}

// 	// каждую задачу перевести в stopped
// 	for _, issue := range issues.Issues {
// 		issueUpdate := IssueUpdateObject{
// 			ProjectID:      issue.Project.Id,
// 			TrackerID:      issue.Tracker.Id,
// 			StatusID:       6,
// 			PriorityID:     issue.Priority.Id,
// 			Subject:        issue.Subject,
// 			Description:    issue.Description,
// 			CategoryID:     issue.Category.Id,
// 			FixedVersionID: issue.FixedVersion.Id,
// 			AssignedToID:   issue.AssignedTo.Id,
// 			ParentIssueID:  issue.Parent.ID,
// 			IsPrivate:      issue.IsPrivate,
// 			EstimatedHours: issue.EstimatedHours,
// 		}
// 		err := c.IssueUpdate(ctx, issue.ID, issueUpdate)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (c *Redmine) IssueToInprogress(ctx context.Context, userId int, issueId int) error {

// 	// каждую задачу перевести в in progress
// 	issueUpdate := IssueUpdateObject{
// 		StatusID: 2,
// 	}

// 	err := c.IssueUpdate(ctx, issueId, issueUpdate)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
