package domain

type BranchIdentifier struct {
	Repository RepositoryIdentifier
	Branch     string
}

func (b *BranchIdentifier) String() string {
	return b.Repository.String() + ":" + b.Branch
}

type Branch struct {
	BranchIdentifier
	Commit CommitIdentifier
}

type CommitIdentifier struct {
	Repository RepositoryIdentifier
	SHA        string
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
	Repository RepositoryIdentifier
	SHA        string
}
