package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func NewLocalGit() (*LocalGit, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("can not open local git: %w", err)
	}

	return &LocalGit{
		repo: r,
	}, nil
}

type LocalGit struct {
	repo *git.Repository
}

func (lg *LocalGit) Remote() (string, error) {
	r, err := lg.repo.Remote("origin")
	if err != nil {
		return "", fmt.Errorf("can not find remote: %w", err)
	}

	c := r.Config()
	if len(c.URLs) < 1 {
		return "", fmt.Errorf("no remote in git repo")
	}

	return c.URLs[0], nil
}

func (lg *LocalGit) CurrentBranch() (string, error) {
	r, err := lg.repo.Head()
	if err != nil {
		return "", fmt.Errorf("can not find current branch: %w", err)
	}

	refName := r.Name()
	return refName.Short(), nil
}
