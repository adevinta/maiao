package credentials

import "fmt"

// Credentials defines the authentication credentials values
type Credentials struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// CredentialGetter defines to implement to get credentials
type CredentialGetter interface {
	CredentialForHost(string) (*Credentials, error)
}

type ChainCredentialGetter []CredentialGetter

func (c ChainCredentialGetter) CredentialForHost(host string) (*Credentials, error) {
	errors := Errors{}
	for _, getter := range c {
		c, err := getter.CredentialForHost(host)
		if err != nil {
			errors = append(errors, err)
		} else {
			return c, nil
		}
	}
	return nil, errors
}

type Errors []error

// Error implements the error interface
func (e Errors) Error() string {
	msg := ""
	sep := ""
	for _, err := range e {
		msg = fmt.Sprintf("%s%s%s", msg, sep, err.Error())
		sep = "\n"
	}
	return msg
}
