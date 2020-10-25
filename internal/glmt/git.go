package glmt

type Git interface {
	Remote() (string, error)
	CurrentBranch() (string, error)
}
