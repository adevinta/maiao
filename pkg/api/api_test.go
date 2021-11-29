package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepoFromGitRemote(t *testing.T) {
	testNewRepoFromGitRemoteResults(t, "https://github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "https://user@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com", Username: "user"})
	testNewRepoFromGitRemoteResults(t, "https://user:password@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com", Username: "user", Password: "password"})
	testNewRepoFromGitRemoteResults(t, "http://github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "http://user@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com", Username: "user"})
	testNewRepoFromGitRemoteResults(t, "http://user:password@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com", Username: "user", Password: "password"})
	testNewRepoFromGitRemoteResults(t, "git@github.company.example.com:maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "git://github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "git://github.company.example.com/maiao/maiao", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "git://user@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "git://user:password@github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
	testNewRepoFromGitRemoteResults(t, "ssh://github.company.example.com/maiao/maiao.git", Repo{Repository: "maiao/maiao", Domain: "github.company.example.com"})
}

func testNewRepoFromGitRemoteResults(t *testing.T, remote string, expected Repo) {
	t.Run(fmt.Sprintf("with a remote of %s", remote), func(t *testing.T) {
		repo, err := NewRepoFromGitRemote(remote)
		assert.NoError(t, err)
		assert.Equal(t, expected, *repo)
	})
}
