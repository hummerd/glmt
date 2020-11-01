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
	URL       string    `json:"url"`
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
