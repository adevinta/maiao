package gh

import (
	"fmt"
	"os"
	"testing"

	"github.com/adevinta/maiao/pkg/credentials"
	"github.com/adevinta/maiao/pkg/system"
	"github.com/stretchr/testify/assert"
)

func setDefaultCredentialStore(c credentials.CredentialGetter) {
	DefaultCredentialGetter = c
}

func setEnterpriseCredentialStore(c credentials.CredentialGetter) {
	GitHubEnterpriseCredentialGetter = c
}

type TestCredentialGetter struct {
	Credentials *credentials.Credentials
	Error       error
	Check       func()
}

type GithubTokenTestData struct {
	setCredentialStore  func(credentials.CredentialGetter)
	domain              string
	environmentVariable string
}

func (c *TestCredentialGetter) CredentialForHost(string) (*credentials.Credentials, error) {
	if c.Check != nil {
		c.Check()
	}
	return c.Credentials, c.Error
}

func sprintf(original string, githubGithubTokenTestData GithubTokenTestData) string {
	return fmt.Sprintf(original+" (domain: '%s', environment variable: '%s'", githubGithubTokenTestData.domain, githubGithubTokenTestData.environmentVariable)
}

func GitHubTokenTest(t *testing.T, githubTokenTestData GithubTokenTestData) {
	t.Cleanup(system.Reset)
	os.Unsetenv(githubTokenTestData.environmentVariable)
	creds := &TestCredentialGetter{}
	defer githubTokenTestData.setCredentialStore(DefaultCredentialGetter)
	githubTokenTestData.setCredentialStore(creds)
	t.Run(sprintf("when username and password are provided, password is used as token", githubTokenTestData), func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Username: "user",
			Password: "api key",
		}
		token, err := getGithubToken(githubTokenTestData.domain)
		assert.NoError(t, err)
		assert.Equal(t, "api key", token)
	})
	t.Run(sprintf("when username only is provided, username is used as token", githubTokenTestData), func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Username: "user",
		}
		token, err := getGithubToken(githubTokenTestData.domain)
		assert.NoError(t, err)
		assert.Equal(t, "user", token)
	})
	t.Run(sprintf("when password only is provided, password is used as token", githubTokenTestData), func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Password: "api key",
		}
		token, err := getGithubToken(githubTokenTestData.domain)
		assert.NoError(t, err)
		assert.Equal(t, "api key", token)
	})
	t.Run(sprintf("when nothing is provided, an error is returned", githubTokenTestData), func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{}
		token, err := getGithubToken(githubTokenTestData.domain)
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})
}

func TestGetGitHubDotComToken(t *testing.T) {
	GitHubTokenTest(t, GithubTokenTestData{
		setCredentialStore:  setDefaultCredentialStore,
		domain:              "github.com",
		environmentVariable: "GITHUB_TOKEN",
	})
}

func TestGetGitHubEnterpriseToken(t *testing.T) {
	GitHubTokenTest(t, GithubTokenTestData{
		setCredentialStore:  setEnterpriseCredentialStore,
		domain:              "test.domain.tld",
		environmentVariable: "GITHUB_ENTERPRISE_TOKEN",
	})
}
