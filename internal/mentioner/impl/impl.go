// Package impl implements services for mentioner.
package impl

import (
	"math/rand"
	"strings"
	"time"

	"gitlab.com/gitlab-merge-tool/glmt/internal/mentioner"
)

type mentionDistributor struct {
	mentionsConfig   *mentionsConfig
	mentionsCount    int
	currentMemberUID string
}

func (d mentionDistributor) distribute(project string) []mentioner.Member {
	if len(d.mentionsConfig.Members) == 0 || d.mentionsCount <= 0 {
		return nil
	}

	members := d.mentionsConfig.Members
	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	owners := make([]mentioner.Member, 0, len(d.mentionsConfig.Members))
	restMembers := make([]mentioner.Member, 0, len(d.mentionsConfig.Members))

	for _, member := range members {
		switch {
		case !member.IsActive:
			continue
		case strings.EqualFold(member.UID, d.currentMemberUID):
			continue
		}

		var isOwner bool
		for _, r := range member.OwnsProjects {
			if strings.EqualFold(r, project) {
				isOwner = true
				break
			}
		}

		switch {
		case isOwner:
			owners = append(owners, member)
		default:
			restMembers = append(restMembers, member)
		}
	}

	selectedMembers := append(owners, restMembers...)
	if len(selectedMembers) > d.mentionsCount {
		selectedMembers = selectedMembers[:d.mentionsCount]
	}

	return selectedMembers
}

type mentionsConfig struct {
	Members []mentioner.Member `json:"members"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
