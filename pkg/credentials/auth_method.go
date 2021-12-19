package credentials

import (
	"net/http"

	"github.com/adevinta/maiao/pkg/log"
)

type GitAuth struct {
	Credentials CredentialGetter
}

func (a *GitAuth) SetAuth(r *http.Request) {
	creds, err := a.Credentials.CredentialForHost(r.Host)
	if err != nil || creds == nil {
		log.Logger.WithField("host", r.Host).WithError(err).Infof("failed to find credentials")
		return
	}
	r.SetBasicAuth(creds.Username, creds.Password)
}

func (a *GitAuth) Name() string {
	return "auth from credentials"
}

func (a *GitAuth) String() string {
	return "auth from credentials"
}
