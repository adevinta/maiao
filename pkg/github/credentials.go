package gh

import (
	"os"

	"github.com/99designs/keyring"
	"github.com/adevinta/maiao/pkg/credentials"
)

// DefaultCredentialGetter implements retrieving credentials from a netrc formatted
// file, with a location of ~/.netrc
var DefaultCredentialGetter credentials.CredentialGetter = defaultCredentials()

func defaultCredentials() credentials.CredentialGetter {
	getters := []credentials.CredentialGetter{
		&credentials.EnvToken{PasswordKey: "GITHUB_TOKEN"},
		&credentials.Netrc{},
		&credentials.GitCredentials{GitPath: "git"},
	}
	if os.Getenv("MAIAO_EXPERIMENTAL_CREDENTIALS") == "true" {
		getters = append(getters,
			credentials.MustNewKeyring(keyring.Config{
				ServiceName:              "maiao",
				PassPrefix:               "maiao/",
				KeychainTrustApplication: true,
				KeychainSynchronizable:   true,
			}),
		)
	}
	return credentials.ChainCredentialGetter(getters)
}
