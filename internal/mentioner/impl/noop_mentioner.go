package impl

import "gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"

func NewNOOPMentioner() *NOOPMentioner {
	return &NOOPMentioner{}
}

type NOOPMentioner struct {
}

func (NOOPMentioner) Mentions(string) ([]mentioner.Member, error) {
	return nil, nil
}
