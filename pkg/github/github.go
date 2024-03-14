package gh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/adevinta/maiao/pkg/log"
	"github.com/google/go-github/v55/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func NewHTTPClientForDomain(ctx context.Context, domain string) (*http.Client, error) {
	if domain == "github.com" {
		domain = "api.github.com"
	}
	// TODO: move this to handle unauthorized calls.
	token, err := getGithubToken(domain)
	if err != nil {
		log.ForContext(ctx).WithError(err).WithField("domain", domain).Errorf("unable to find token")
		return nil, fmt.Errorf("unable to find token for %s: %s", domain, err.Error())
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc, nil
}

// NewClient instanciates a new github client depending on the domain name
//
// When requesting a client for a different host than github.com,
// a client for github enterprise is considered
//
// Credentials are even taken from GITHUB_TOKEN environment variable or
// from your ~/.netrc file
func NewClient(httpClient *http.Client, domain string) (*github.Client, error) {
	c := github.NewClient(httpClient)
	switch domain {
	case "api.github.com", "github.com":
	default:
		GitHubURL := url.URL{
			Scheme: "https",
			Host:   domain,
			Path:   "/api/v3/",
		}
		GitHubUploadURL := GitHubURL
		GitHubUploadURL.Path = "/api/v3/upload/"
		c.BaseURL = &GitHubURL
		// TODO: confirm from https://github.com/goreleaser/goreleaser/issues/365#issuecomment-331655225
		c.UploadURL = &GitHubUploadURL
	}
	return c, nil
}

func getGithubToken(domain string) (string, error) {
	creds, err := DefaultCredentialGetter.CredentialForHost(domain)
	if err != nil {
		return "", err
	}
	token := findFirstNonEmptyString(creds.Password, creds.Username)
	if token != "" {
		Logger.WithFields(logrus.Fields{
			"context": "initializing GitHub client",
		}).Debugf("using github token from credentials store")
		return token, nil
	}
	return "", fmt.Errorf("unable to find a token for domain %s", domain)
}

func findFirstNonEmptyString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
