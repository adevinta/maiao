package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/google/go-github/v40/github"
	"github.com/sirupsen/logrus"
	gh "github.com/adevinta/maiao/pkg/github"
	"github.com/adevinta/maiao/pkg/log"
)

// GitHub implements the PullRequester interface allowing to create pull requests for a given repository
type GitHub struct {
	*github.Client
	Host       string
	Owner      string
	Repository string
}

// Ensure ensures a PR is opened for the head branch
func (g *GitHub) Ensure(ctx context.Context, options PullRequestOptions) (*PullRequest, bool, error) {
	ctx = log.WithContextFields(ctx, logrus.Fields{
		"context":    "ensuring existing pull request",
		"owner":      g.Owner,
		"repository": g.Repository,
		"prOptions":  options,
	})
	prs, _, err := g.PullRequests.List(context.Background(), g.Owner, g.Repository, &github.PullRequestListOptions{
		Head:      g.Owner + ":" + options.Head,
		Sort:      "created",
		Direction: "desc",
	})
	if err != nil {
		log.ForContext(ctx).WithError(err).Error("failed to list existing pull requests")
		return nil, false, err
	}
	switch len(prs) {
	case 0:
		pr, _, err := g.PullRequests.Create(context.Background(), g.Owner, g.Repository, &github.NewPullRequest{
			Title: github.String(options.Title),
			Body:  github.String(options.Body),
			Base:  github.String(options.Base),
			Head:  github.String(options.Head),
		})
		if err != nil {
			log.ForContext(ctx).WithError(err).Error("failed to create new pull request")
			return nil, false, err
		}
		log.ForContext(ctx).Debug("new PR has been created")
		return &PullRequest{
			ID:  strconv.Itoa(*pr.Number),
			URL: pr.GetHTMLURL(),
		}, true, nil
	case 1:
		log.ForContext(ctx).Trace("PR already existed")
		return &PullRequest{
			ID:  strconv.Itoa(*prs[0].Number),
			URL: prs[0].GetHTMLURL(),
		}, false, nil
	}
	log.ForContext(ctx).WithError(err).Error("failed to list existing pull requests")
	return nil, false, errors.New("Too may matching pull requests")

}

// Update implements the Update interface to update an existing pull request
func (g *GitHub) Update(ctx context.Context, pr *PullRequest, options PullRequestOptions) (*PullRequest, error) {
	ctx = log.WithContextFields(ctx, logrus.Fields{
		"context":    "ensuring existing pull request",
		"owner":      g.Owner,
		"repository": g.Repository,
		"prOptions":  options,
	})
	id, err := strconv.Atoi(pr.ID)
	if err != nil {
		log.ForContext(ctx).WithField("prID", pr.ID).WithError(err).Error("failed to parse pull request ID")
		return nil, err
	}
	ctx = log.WithContextFields(ctx, logrus.Fields{"prID": id})
	p, _, err := g.PullRequests.Edit(ctx, g.Owner, g.Repository, id, &github.PullRequest{
		Title: github.String(options.Title),
		Body:  github.String(options.Body),
		Base: &github.PullRequestBranch{
			Ref: github.String(options.Base),
		},
		Head: &github.PullRequestBranch{
			Ref: github.String(options.Head),
		},
	})
	if err != nil {
		log.ForContext(ctx).WithError(err).Error("failed to edit pull request")
		return nil, err
	}
	log.ForContext(ctx).Info("edit pull request")
	return &PullRequest{
		ID:  strconv.Itoa(*p.Number),
		URL: *p.URL,
	}, err
}

// LinkedTopicIssues returns the search URL for linked issues
func (g *GitHub) LinkedTopicIssues(topic string) string {
	return `https://` + g.Host + `/search?q=type%3Apr+%22Topic%3A+` + url.QueryEscape(topic) + `%22&type=Issues`
}

// NewGitHubUpserter instanciates an upserter that uses the github API to create and update pull requests
func NewGitHubUpserter(ctx context.Context, endpoint *transport.Endpoint) (*GitHub, error) {
	ctx = log.WithContextFields(ctx, logrus.Fields{
		"context":  "initializing GitHub client",
		"endpoint": endpoint,
	})

	orgRepo := strings.Split(strings.Trim(endpoint.Path, "/"), "/")
	if len(orgRepo) != 2 {
		log.ForContext(ctx).WithField("repository", endpoint.Path).Error("invalid repository, expecting <org>/<repo>")
		return nil, fmt.Errorf("invalid repository, expecting <org>/<repo>")
	}
	client, err := gh.NewClient(endpoint.Host)
	if err != nil {
		log.ForContext(ctx).WithError(err).Errorf("failed to create a new github client: %s", err.Error())
		return nil, err
	}
	repo, _, err := client.Repositories.Get(ctx, orgRepo[0], strings.TrimSuffix(orgRepo[1], ".git"))
	if err != nil {
		return nil, err
	}

	gh := &GitHub{
		Host:       endpoint.Host,
		Owner:      repo.GetOwner().GetLogin(),
		Repository: repo.GetName(),
		Client:     client,
	}
	log.ForContext(ctx).Trace("initialized github client")
	return gh, nil
}
