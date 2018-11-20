package domain

type BranchIdentifier struct {
	RepositoryIdentifier
	Branch string
}

type Branch struct {
	BranchIdentifier
	Commit CommitIdentifier
}

type CommitIdentifier struct {
	RepositoryIdentifier
	SHA string
}

type Commit struct {
	CommitIdentifier
	Message string
	Parents []CommitIdentifier
	Tree    TreeIdentifier
}

func (c *Commit) GetSingleParent() *CommitIdentifier {
	return nil
}

type TreeIdentifier struct {
	RepositoryIdentifier
	SHA string
}
