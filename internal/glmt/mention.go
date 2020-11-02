package glmt

import (
	"math/rand"
	"strings"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

// Mentions selects users to be mentioned in MR. Sekects project owner (if there is one)
// plus additional random members.
func Mentions(t *team.Team, me, project string, count int) []*team.Member {
	if t == nil {
		return nil
	}

	members := t.Members
	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	owners := make([]*team.Member, 0, len(members))
	restMembers := make([]*team.Member, 0, len(members))

	me = strings.TrimPrefix(me, "@")

	for _, member := range members {
		switch {
		case !member.IsActive:
			continue
		case strings.EqualFold(member.Username, me):
			continue
		}

		var isOwner bool
		for _, r := range member.OwnsProjects {
			if strings.EqualFold(r, project) {
				isOwner = true
				break
			}
		}

		if isOwner {
			owners = append(owners, member)
		} else {
			restMembers = append(restMembers, member)
		}
	}

	selectedMembers := append(owners, restMembers...)
	if len(selectedMembers) > count {
		selectedMembers = selectedMembers[:count]
	}

	return selectedMembers
}
