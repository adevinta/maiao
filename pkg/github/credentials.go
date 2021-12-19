package gh

import (
	"github.com/adevinta/maiao/pkg/credentials"
)

// DefaultCredentialGetter implements retrieving credentials from a netrc formatted
// file, with a location of ~/.netrc
var DefaultCredentialGetter credentials.CredentialGetter = credentials.ChainCredentialGetter([]credentials.CredentialGetter{
	&credentials.EnvToken{PasswordKey: "GITHUB_TOKEN"},
	&credentials.Netrc{},
	&credentials.GitCredentials{GitPath: "git"},
})
