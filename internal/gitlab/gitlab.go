// Package gitlab defines interface for GitLab service
package gitlab

import (
	"context"
	"time"
)

type CreateMRRequest struct {
	Project            string `json:"id"`
	SourceBranch       string `json:"source_branch"`
	TargetBranch       string `json:"target_branch"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Squash             bool   `json:"squash"`
	RemoveSourceBranch bool   `json:"remove_source_branch"`
}

type CreateMRResponse struct {
	ID        int64     `json:"id"`
	IID       int64     `json:"iid"`
	ProjectID int64     `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"web_url"`
	// ChangesCount shows how many changes in MR
	// Note: the value in the response is a string, not an integer. This is because when an MR
	// has too many changes to display and store, it will be capped at 1,000. In that case,
	// the API will return the string "1000+" for the changes count.
	ChangesCount string `json:"changes_count"`
}

type GitlabError struct {
	Message string
}

func (e GitlabError) Error() string {
	return e.Message
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type GitLab interface {
	CreateMR(ctx context.Context, req CreateMRRequest) (CreateMRResponse, error)
	CurrentUser(ctx context.Context) (UserResponse, error)
}
