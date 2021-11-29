package credentials

import (
	"fmt"
	"os/user"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/adevinta/maiao/pkg/system"
)

func testCredentialSuccess(t *testing.T, n *Netrc, machine string, expected Credentials) {
	fmt.Println(n.Path)
	c, err := n.CredentialForHost(machine)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, expected, *c)
}

func TestNetrcCredentials(t *testing.T) {
	var n *Netrc
	t.Run("when netrc handler is null, an error is returned", func(t *testing.T) {
		c, err := n.CredentialForHost("some host")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
	n = &Netrc{}
	t.Run("when failing to get current user, an error is returned", func(t *testing.T) {
		t.Cleanup(system.Reset)
		system.CurrentUser = func() (*user.User, error) {
			return nil, fmt.Errorf("test error")
		}
		c, err := n.CredentialForHost("some host")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
	n.CredentialForHost("some host")
	t.Run("when no path is provided, a default one is created", func(t *testing.T) {
		if !strings.HasSuffix(n.Path, "/.netrc") {
			t.Errorf("default path %s does not have the /.netrc suffix", n.Path)
		}
	})
	n.Path = "resources/.netrc"
	t.Run("when a path is provided", func(t *testing.T) {
		t.Run("and the machine has only login, credentials is returned", func(t *testing.T) {
			testCredentialSuccess(t, n, "login.example.com", Credentials{Username: "this-is-a-login"})
		})
		t.Run("and the machine has only password, credentials is returned", func(t *testing.T) {
			testCredentialSuccess(t, n, "password.example.com", Credentials{Password: "pass"})
		})
		t.Run("and the machine has login and password, credentials is returned", func(t *testing.T) {
			testCredentialSuccess(t, n, "example.com", Credentials{Username: "me", Password: "a-secure-password"})
		})
		t.Run("when the path does not exist, an eror is returned", func(t *testing.T) {
			n.Path = "unavailable file"
			c, err := n.CredentialForHost("any")
			assert.Error(t, err)
			assert.Nil(t, c)
		})
	})
}
