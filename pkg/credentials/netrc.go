package credentials

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jdxcode/netrc"
	"github.com/sirupsen/logrus"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/system"
)

// Netrc implements the CredentialGetter interface,
// getting the credentials from a netrc formatted file.
// When path is empty, the default ~/.netrc path is used
type Netrc struct {
	Path string
}

// CredentialForHost retrieves the credentials for a given host in the netrc file
func (n *Netrc) CredentialForHost(host string) (*Credentials, error) {
	if n == nil {
		return nil, fmt.Errorf("failed to find credentials for machine %s, nil handler", host)
	}
	ctx := log.WithContextFields(context.Background(), logrus.Fields{"context": "parsing netrc", "host": host})
	logger := log.ForContext(ctx)
	if n.Path == "" {
		usr, err := system.CurrentUser()
		if err != nil {
			logger.WithError(err).Infof("failed to retrieve current user")
			return nil, err
		}
		n.Path = filepath.Join(usr.HomeDir, ".netrc")
		logger.Debugf("using default netrc path")
	}
	logger = logger.WithContext(log.WithContextFields(ctx, logrus.Fields{"path": n.Path}))
	parsed, err := netrc.Parse(n.Path)
	if err != nil {
		logger.WithError(err).Infof("failed to parse netrc file")
		return nil, err
	}
	machine := parsed.Machine(host)
	if machine != nil {
		logger.Debugf("found credentials")
		return &Credentials{
			Username: machine.Get("login"),
			Password: machine.Get("password"),
		}, nil
	}
	return nil, fmt.Errorf("failed to find credentials for host %s in netRC %s", host, n.Path)
}
