package maiao

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/adevinta/maiao/pkg/log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRepository struct {
	remote          func(name string) (*git.Remote, error)
	branches        func() (storer.ReferenceIter, error)
	config          func() (*config.Config, error)
	log             func(o *git.LogOptions) (object.CommitIter, error)
	resolveRevision func(rev plumbing.Revision) (*plumbing.Hash, error)
	fetch           func(o *git.FetchOptions) error
	push            func(o *git.PushOptions) error
	worktree        func() (*git.Worktree, error)
}

func (r *testRepository) Head() (*plumbing.Reference, error) {
	return nil, errors.New("not implemented")
}

func (r *testRepository) Remote(name string) (*git.Remote, error) {
	if r.config == nil {
		return nil, errors.New("not implemented")
	}
	return r.remote(name)
}

func (r *testRepository) Push(o *git.PushOptions) error {
	if r.config == nil {
		return errors.New("not implemented")
	}
	return r.push(o)
}

func (r *testRepository) Branches() (storer.ReferenceIter, error) {
	if r.config == nil {
		return nil, errors.New("not implemented")
	}
	return r.branches()
}

func (r *testRepository) Config() (*config.Config, error) {
	if r.config == nil {
		return nil, errors.New("not implemented")
	}
	return r.config()
}

func (r *testRepository) Fetch(o *git.FetchOptions) error {
	if r.fetch == nil {
		return errors.New("not implemented")
	}
	return r.fetch(o)
}

func (r *testRepository) Log(o *git.LogOptions) (object.CommitIter, error) {
	if r.log == nil {
		return nil, errors.New("not implemented")
	}
	return r.log(o)
}

func (r *testRepository) ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error) {
	if r.resolveRevision == nil {
		return nil, errors.New("not implemented")
	}
	return r.resolveRevision(rev)
}

func (r *testRepository) Worktree() (*git.Worktree, error) {
	if r.worktree == nil {
		return nil, errors.New("not implemented")
	}
	return r.worktree()
}

func TestDefaultOptionsUsesGitDefaults(t *testing.T) {
	opts := ReviewOptions{}
	defaultOptions(context.Background(), &testRepository{}, &opts)
	assert.Equal(t, "master", opts.Branch)
	assert.Equal(t, "origin", opts.Remote)
}

func TestDefaultOptionsUsesTrackingRemote(t *testing.T) {
	branch := uuid.New().String()
	remoteName := uuid.New().String()
	opts := ReviewOptions{
		Branch: branch,
	}
	defaultOptions(context.Background(), &testRepository{
		config: func() (*config.Config, error) {
			return &config.Config{
				Branches: map[string]*config.Branch{
					branch: {Remote: remoteName},
				},
			}, nil
		},
	}, &opts)
	assert.Equal(t, branch, opts.Branch)
	assert.Equal(t, remoteName, opts.Remote)
}

func TestDefaultOptionsUsesGitDefaultRemoteNameWhenNotTracked(t *testing.T) {
	branch := uuid.New().String()
	opts := ReviewOptions{
		Branch: branch,
	}
	defaultOptions(context.Background(), &testRepository{
		config: func() (*config.Config, error) {
			return &config.Config{
				Branches: map[string]*config.Branch{
					branch: {},
				},
			}, nil
		},
	}, &opts)
	assert.Equal(t, branch, opts.Branch)
	assert.Equal(t, "origin", opts.Remote)
}

func TestDefaultOptionsUsesGitDefaultRemoteNameWhenTheBranchIsNotFound(t *testing.T) {
	branch := uuid.New().String()
	opts := ReviewOptions{
		Branch: branch,
	}
	defaultOptions(context.Background(), &testRepository{
		config: func() (*config.Config, error) {
			return &config.Config{
				Branches: map[string]*config.Branch{},
			}, nil
		},
	}, &opts)
	assert.Equal(t, branch, opts.Branch)
	assert.Equal(t, "origin", opts.Remote)
}

func TestDefaultOptionsUsesGitDefaultRemoteNameWhenTheBranchConfigIsNull(t *testing.T) {
	branch := uuid.New().String()
	opts := ReviewOptions{
		Branch: branch,
	}
	defaultOptions(context.Background(), &testRepository{
		config: func() (*config.Config, error) {
			return &config.Config{
				Branches: map[string]*config.Branch{
					branch: nil,
				},
			}, nil
		},
	}, &opts)
	assert.Equal(t, branch, opts.Branch)
	assert.Equal(t, "origin", opts.Remote)
}

func TestRebaseTODO(t *testing.T) {
	assert.Equal(
		t,
		strings.Join([]string{}, "\n"),
		rebaseTODO([]*change{}),
	)
	assert.Equal(
		t,
		strings.Join([]string{
			"reword b34ccd81a342e155b8382992cddb116c56bee95c other-change",
			"pick c30a2f070b4f3d00c26679186345ea506e664056 fixup! other-change",
			"pick 943c8d8469c2800e361cea0f37a3e38cc7e90fd6 add hello world",
		}, "\n"),
		rebaseTODO([]*change{
			{
				commits: []*object.Commit{
					{Hash: plumbing.NewHash("b34ccd81a342e155b8382992cddb116c56bee95c"), Message: "other-change"},
					{Hash: plumbing.NewHash("c30a2f070b4f3d00c26679186345ea506e664056"), Message: "fixup! other-change"},
				},
			},
			{
				changeID: "1234",
				commits: []*object.Commit{
					{Hash: plumbing.NewHash("943c8d8469c2800e361cea0f37a3e38cc7e90fd6"), Message: "add hello world"},
				},
			},
			{},
		}),
	)
}

func TestRemoveMergedChangeIDs(t *testing.T) {
	assert.Equal(
		t,
		[]*change{
			{
				changeID: "5678",
				commits: []*object.Commit{
					{Hash: plumbing.NewHash("b34ccd81a342e155b8382992cddb116c56bee95c"), Message: "other-change"},
					{Hash: plumbing.NewHash("c30a2f070b4f3d00c26679186345ea506e664056"), Message: "fixup! other-change"},
				},
			},
			{
				commits: []*object.Commit{
					{Hash: plumbing.NewHash("943c8d8469c2800e361cea0f37a3e38cc7e90fd6"), Message: "add hello world"},
				},
			},
			{},
		},
		removeMergedChangeIDs(
			[]*change{
				{
					changeID: "1234",
					commits: []*object.Commit{
						{Hash: plumbing.NewHash("5b380f1b4081a7b64b72954a6ad58f12131749ed"), Message: "merged-change"},
					},
				},
				{
					changeID: "5678",
					commits: []*object.Commit{
						{Hash: plumbing.NewHash("b34ccd81a342e155b8382992cddb116c56bee95c"), Message: "other-change"},
						{Hash: plumbing.NewHash("c30a2f070b4f3d00c26679186345ea506e664056"), Message: "fixup! other-change"},
					},
				},
				{
					commits: []*object.Commit{
						{Hash: plumbing.NewHash("943c8d8469c2800e361cea0f37a3e38cc7e90fd6"), Message: "add hello world"},
					},
				},
				{},
			},
			map[string]struct{}{
				"1234": {},
			},
		),
	)
}

func TestNeedReview(t *testing.T) {
	storage := memory.NewStorage()
	rootParent := "bdc945b1bc57b3938f7223c7adb8bc2db58b838f"
	newTestCommit(t, storage, rootParent)
	assert.False(t, changesNeedRebase(context.Background(), []*change{}))
	t.Run("When change ID is missing", func(t *testing.T) {
		assert.True(t, changesNeedRebase(context.Background(), []*change{
			{},
		}))
	})
	t.Run("When change ID exists", func(t *testing.T) {
		assert.False(t, changesNeedRebase(context.Background(), []*change{
			{
				changeID: "changeID",
			},
		}))
	})
	t.Run("When fixups are in order", func(t *testing.T) {
		assert.False(t, changesNeedRebase(context.Background(), []*change{
			{
				changeID: "changeID",
				commits: []*object.Commit{
					newTestCommit(t, storage, "fc73d3a47b5864a8668eb826d506deb2bb54c1b5", rootParent),
					newTestCommit(t, storage, "2d885b9b60dd70bb5c9b66ac72d21da894787fd7", "fc73d3a47b5864a8668eb826d506deb2bb54c1b5"),
				},
			},
			{
				changeID: "changeID2",
				commits: []*object.Commit{
					newTestCommit(t, storage, "832e80c48ef021d64a6d38450584f0e7f6f333a2", "2d885b9b60dd70bb5c9b66ac72d21da894787fd7"),
					newTestCommit(t, storage, "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d", "832e80c48ef021d64a6d38450584f0e7f6f333a2"),
				},
			},
		}))
	})
	t.Run("When fixups are out of order", func(t *testing.T) {
		assert.True(t, changesNeedRebase(context.Background(), []*change{
			{
				changeID: "changeID",
				commits: []*object.Commit{
					newTestCommit(t, storage, "fc73d3a47b5864a8668eb826d506deb2bb54c1b5", rootParent),
					newTestCommit(t, storage, "2d885b9b60dd70bb5c9b66ac72d21da894787fd7", "fc73d3a47b5864a8668eb826d506deb2bb54c1b5"),
					newTestCommit(t, storage, "61809666bdb341715a6df144a1f84a6709975816", "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d"),
				},
			},
			{
				changeID: "changeID2",
				commits: []*object.Commit{
					newTestCommit(t, storage, "832e80c48ef021d64a6d38450584f0e7f6f333a2", "2d885b9b60dd70bb5c9b66ac72d21da894787fd7"),
					newTestCommit(t, storage, "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d", "832e80c48ef021d64a6d38450584f0e7f6f333a2"),
				},
			},
		}))
	})
	t.Run("When main commits are out of order", func(t *testing.T) {
		assert.True(t, changesNeedRebase(context.Background(), []*change{
			{
				changeID: "changeID",
				commits: []*object.Commit{
					newTestCommit(t, storage, "fc73d3a47b5864a8668eb826d506deb2bb54c1b5", "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d"),
					newTestCommit(t, storage, "2d885b9b60dd70bb5c9b66ac72d21da894787fd7", "fc73d3a47b5864a8668eb826d506deb2bb54c1b5"),
					newTestCommit(t, storage, "61809666bdb341715a6df144a1f84a6709975816", "2d885b9b60dd70bb5c9b66ac72d21da894787fd7"),
				},
			},
			{
				changeID: "changeID2",
				commits: []*object.Commit{
					newTestCommit(t, storage, "832e80c48ef021d64a6d38450584f0e7f6f333a2", rootParent),
					newTestCommit(t, storage, "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d", "832e80c48ef021d64a6d38450584f0e7f6f333a2"),
				},
			},
		}))
	})
	t.Run("When one parent is not found", func(t *testing.T) {
		storage := memory.NewStorage()
		assert.True(t, changesNeedRebase(context.Background(), []*change{
			{
				changeID: "changeID",
				commits: []*object.Commit{
					newTestCommit(t, storage, "832e80c48ef021d64a6d38450584f0e7f6f333a2", rootParent),
					newTestCommit(t, storage, "fc73d3a47b5864a8668eb826d506deb2bb54c1b5", "b1b5506d096e9697ca5ad8fee28ce0415ff4bc0d"),
				},
			},
		}))
	})
}

func newTestCommit(t *testing.T, storer storer.EncodedObjectStorer, sha string, parents ...string) *object.Commit {
	t.Helper()
	c := *&object.Commit{
		Hash: plumbing.NewHash(sha),
	}
	for _, parent := range parents {
		c.ParentHashes = append(c.ParentHashes, plumbing.NewHash(parent))
	}
	o := &testEncodedObject{
		sha:  sha,
		data: &bytes.Buffer{},
	}
	require.NoError(t, c.EncodeWithoutSignature(o))
	storer.SetEncodedObject(o)
	out, err := object.GetCommit(storer, plumbing.NewHash(sha))
	require.NoError(t, err)
	return out
}

type testEncodedObject struct {
	sha        string
	objectType plumbing.ObjectType
	data       *bytes.Buffer
	size       int64
}

func (t *testEncodedObject) Hash() plumbing.Hash {
	return plumbing.NewHash(t.sha)
}
func (t *testEncodedObject) Type() plumbing.ObjectType {
	return t.objectType
}
func (t *testEncodedObject) SetType(objectType plumbing.ObjectType) {
	t.objectType = objectType
}
func (t *testEncodedObject) Size() int64 {
	return t.size
}
func (t *testEncodedObject) SetSize(size int64) {
	t.size = size
}
func (t *testEncodedObject) Reader() (io.ReadCloser, error) {
	return ioutil.NopCloser(t.data), nil
}
func (t *testEncodedObject) Writer() (io.WriteCloser, error) {
	return &nopWriteCloser{Writer: t.data}, nil
}

type nopWriteCloser struct {
	io.Writer
}

func (n *nopWriteCloser) Close() error {
	return nil
}

func init() {
	log.Logger.SetLevel(logrus.TraceLevel)
}
