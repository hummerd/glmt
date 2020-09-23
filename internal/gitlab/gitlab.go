// Package gitlab defines interface for GitLab service
package gitlab

import "context"

type CreateMRRequest struct {
	ID                 string `json:"id"`
	SourceBranch       string `json:"source_branch"`
	TargetBranch       string `json:"target_branch"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Squash             bool   `json:"squash"`
	RemoveSourceBranch bool   `json:"remove_source_branch"`
}

type GitLab interface {
	CreateMR(ctx context.Context, req CreateMRRequest) (int64, error)
}
