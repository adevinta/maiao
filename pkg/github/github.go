package gh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v40/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// NewClient instanciates a new github client depending on the domain name
//
// When requesting a client for a different host than github.com,
// a client for github enterprise is considered
//
// Credentials are even taken from GITHUB_TOKEN environment variable or
// from your ~/.netrc file
func NewClient(domain string) (*github.Client, error) {
	logger := Logger.WithFields(logrus.Fields{
		"context": "initializing GitHub client",
	})
	if domain == "github.com" {
		domain = "api.github.com"
	}
	// TODO: move this to handle unauthorized calls.
	token, err := getGithubToken(domain)
	if err != nil {
		logger.Errorf("unable to find token for %s: %s", domain, err.Error())
		return nil, fmt.Errorf("unable to find token for %s: %s", domain, err.Error())
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)
	if domain != "api.github.com" {
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
