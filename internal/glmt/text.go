package glmt

const (
	TmpVarProjectName      = "ProjectName"
	TmpVarBranchName       = "BranchName"
	TmpVarTargetBranchName = "TargetBranchName"
	TmpVarTitle            = "Title"
	TmpVarDescription      = "Description"
	TmpVarMRURL            = "MergeRequestURL"
)

func getTextArgs(branch, projectName string, params CreateMRParams) map[string]string {
	r := map[string]string{}

	defer func() {
		// in the end override values with well known
		r[TmpVarProjectName] = projectName
		r[TmpVarBranchName] = branch
		r[TmpVarTargetBranchName] = params.TargetBranch
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
