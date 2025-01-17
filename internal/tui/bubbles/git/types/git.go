package types

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo interface {
	Name() string
	GetHEAD() *plumbing.Reference
	SetHEAD(*plumbing.Reference) error
	GetReferences() []*plumbing.Reference
	GetReadme() string
	GetReadmePath() string
	GetCommits(*plumbing.Reference) (Commits, error)
	Repository() *git.Repository
	Tree(*plumbing.Reference, string) (*object.Tree, error)
	PatchCtx(context.Context, *object.Commit) (*object.Patch, error)
}

type Commits []*object.Commit

func (cl Commits) Len() int      { return len(cl) }
func (cl Commits) Swap(i, j int) { cl[i], cl[j] = cl[j], cl[i] }
func (cl Commits) Less(i, j int) bool {
	return cl[i].Author.When.After(cl[j].Author.When)
}
