// Package mentioner defines interface for mentioners.
package mentioner

// Mentioner selects random members.
type Mentioner interface {
	// Mentions returns N random members including project owners and
	// excluding current member.
	Mentions(project string) (members []Member, err error)
}

// Member to mention.
type Member struct {
	UID            string   `json:"uid"`
	FullName       string   `json:"full_name"`
	SlackUsername  string   `json:"slack_username"`
	GitlabUsername string   `json:"gitlab_username"`
	OwnsProjects   []string `json:"owns_projects"`
	IsActive       bool     `json:"is_active"`
}
