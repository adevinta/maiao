package credentials

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCredentialsReturnsNoCredentialsWhenGitCommandFails(t *testing.T) {
	t.Cleanup(func() {
		run = realRun
	})
	run = func(opts runOpts) (string, error) {
		return "", errors.New("test error")
	}
	gitCreds := GitCredentials{GitPath: "git"}
	creds, err := gitCreds.CredentialForHost("test-host")
	assert.Nil(t, creds)
	assert.Error(t, err)
}

func TestGitCredentialsReturnsNoCredentialsWhenNoGitHelperIsProvided(t *testing.T) {
	t.Cleanup(func() {
		run = realRun
	})
	run = func(opts runOpts) (string, error) {
		return "", nil
	}
	gitCreds := GitCredentials{GitPath: "git"}
	creds, err := gitCreds.CredentialForHost("test-host")
	assert.Nil(t, creds)
	assert.Error(t, err)
}

func TestGitCredentialsReturnsCredentialsWhenCredentialsAreFound(t *testing.T) {
	t.Cleanup(func() {
		run = realRun
	})
	run = func(opts runOpts) (string, error) {
		assert.Equal(t, "git", opts.path)
		assert.Equal(t, []string{"credential", "fill"}, opts.args)
		assert.Equal(t, "protocol=https\nhost=test-host", opts.stdin)
		return "protocol=https\nhost=test-host\nusername=PersonalAccessToken\npassword=secure-password", nil
	}
	gitCreds := GitCredentials{GitPath: "git"}
	creds, err := gitCreds.CredentialForHost("test-host")
	assert.Equal(t, &Credentials{Username: "PersonalAccessToken", Password: "secure-password"}, creds)
	assert.NoError(t, err)
}

func TestGitCredentialsReturnsNoCredentialsWhenCredentialsReturnNoPassowrd(t *testing.T) {
	t.Cleanup(func() {
		run = realRun
	})
	run = func(opts runOpts) (string, error) {
		assert.Equal(t, "git", opts.path)
		assert.Equal(t, []string{"credential", "fill"}, opts.args)
		assert.Equal(t, "protocol=https\nhost=test-host", opts.stdin)
		return "protocol=https\nhost=test-host\nusername=PersonalAccessToken\n", nil
	}
	gitCreds := GitCredentials{GitPath: "git"}
	creds, err := gitCreds.CredentialForHost("test-host")
	assert.Nil(t, creds)
	assert.Error(t, err)
}
