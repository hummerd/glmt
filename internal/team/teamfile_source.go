// Package teamfile describes gitlab team members
package team

import "context"

type Team struct {
	Members []*Member `json:"members"`
}

// Member is a gitlab team member
type Member struct {
	// Username is gitlab user name
	Username     string   `json:"username"`
	OwnsProjects []string `json:"owns_projects"`
	IsActive     bool     `json:"is_active"`
}

type TeamFileSource interface {
	Team(ctx context.Context) (*Team, error)
}
