package glmt

import (
	"strings"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

const (
	TmpVarProjectName          = "ProjectName"
	TmpVarBranchName           = "BranchName"
	TmpVarRemote               = "Remote"
	TmpVarTargetBranchName     = "TargetBranchName"
	TmpVarTitle                = "Title"
	TmpVarDescription          = "Description"
	TmpVarMRURL                = "MergeRequestURL"
	TmpVarGitlabMentions       = "GitlabMentions"
	TmpVarNotificationMentions = "NotificationMentions"
	TmpVarMRChangesCount       = "ChangesCount"
	TmpVarUsername             = "Username"
)

func getTextArgs(branch, projectName, remote, username string, params CreateMRParams, members []*team.Member) map[string]string {
	r := map[string]string{}

	gitlabMentions := make([]string, 0, len(members))
	for _, m := range members {
		gitlabMentions = append(gitlabMentions, "@"+m.Username)
	}

	defer func() {
		// in the end override values with well known
		r[TmpVarProjectName] = projectName
		r[TmpVarBranchName] = branch
		r[TmpVarTargetBranchName] = params.TargetBranch
		r[TmpVarGitlabMentions] = strings.Join(gitlabMentions, ", ")
		r[TmpVarRemote] = remote
		r[TmpVarUsername] = username
	}()

	if params.BranchRegexp == nil {
		return r
	}

	subNames := params.BranchRegexp.SubexpNames()
	if len(subNames) <= 1 {
		return r
	}

	match := params.BranchRegexp.FindStringSubmatch(branch)
	for i := 1; i < len(subNames); i++ {
		m := ""
		if len(match) > i {
			m = match[i]
		}
		r[subNames[i]] = m
	}

	return r
}
