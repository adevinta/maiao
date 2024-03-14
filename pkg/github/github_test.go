package gh

import (
	"context"
	"os"
	"testing"

	"github.com/adevinta/maiao/pkg/credentials"
	"github.com/adevinta/maiao/pkg/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setCredentialStore(c credentials.CredentialGetter) {
	DefaultCredentialGetter = c
}

type TestCredentialGetter struct {
	Credentials *credentials.Credentials
	Error       error
	Check       func()
}

func TestNewGraphQLClient(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test")
	ctx := context.Background()
	httpClient, err := NewHTTPClientForDomain(ctx, "github.com")
	require.NoError(t, err)
	graphQLClient, err := NewGraphQLClient(httpClient, "github.com")

	assert.NoError(t, err)
	assert.NotNil(t, graphQLClient)

}

func (c *TestCredentialGetter) CredentialForHost(string) (*credentials.Credentials, error) {
	if c.Check != nil {
		c.Check()
	}
	return c.Credentials, c.Error
}

func TestGetGitHubToken(t *testing.T) {
	t.Cleanup(system.Reset)
	os.Unsetenv("GITHUB_TOKEN")
	creds := &TestCredentialGetter{}
	defer setCredentialStore(DefaultCredentialGetter)
	setCredentialStore(creds)
	t.Run("when username and password are provided, password is used as token", func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Username: "user",
			Password: "api key",
		}
		token, err := getGithubToken("test.domain.tld")
		assert.NoError(t, err)
		assert.Equal(t, "api key", token)
	})
	t.Run("when username only is provided, username is used as token", func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Username: "user",
		}
		token, err := getGithubToken("test.domain.tld")
		assert.NoError(t, err)
		assert.Equal(t, "user", token)
	})
	t.Run("when password only is provided, password is used as token", func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{
			Password: "api key",
		}
		token, err := getGithubToken("test.domain.tld")
		assert.NoError(t, err)
		assert.Equal(t, "api key", token)
	})
	t.Run("when nothing is provided, an error is returned", func(t *testing.T) {
		defer func(c *credentials.Credentials) { creds.Credentials = c }(creds.Credentials)
		creds.Credentials = &credentials.Credentials{}
		token, err := getGithubToken("test.domain.tld")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})
}
