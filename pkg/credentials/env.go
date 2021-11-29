package credentials

import (
	"errors"
	"fmt"
	"os"
)

// EnvToken resolves github authentication using the GITHUB_TOKEN environment variable if present
type EnvToken struct {
	UsernameKey     string
	PasswordKey     string
	DefaultUserName string
}

func (e *EnvToken) CredentialForHost(host string) (*Credentials, error) {
	if e.PasswordKey == "" {
		return nil, errors.New("no environment variable to ")
	}
	username := "x-token"
	if e.DefaultUserName != "" {
		username = e.DefaultUserName
	}
	if e.UsernameKey != "" {
		v, ok := os.LookupEnv(e.UsernameKey)
		if ok {
			username = v
		}
	}
	val, ok := os.LookupEnv(e.PasswordKey)
	if ok && val != "" {
		return &Credentials{
			Username: username,
			Password: val,
		}, nil
	}
	return nil, fmt.Errorf("no token found in environment variable %s", e.PasswordKey)
}
