// Package git defines interface for git service
package git

type Git interface {
	Remote() (string, error)
	CurrentBranch() (string, error)
}
