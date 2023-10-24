package maiao

import (
	"testing"

	"github.com/adevinta/maiao/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestDetailsSkipsMissingSummaries(t *testing.T) {
	assert.Equal(t,
		[]string{
			"<details>",
			"hello world",
			"</details>",
		},
		details(
			[]string{
				"hello world",
			},
			"",
		),
	)
}

func TestDetailsIncludesSummary(t *testing.T) {
	assert.Equal(t,
		[]string{
			"<details>",
			"<summary>",
			"summary",
			"</summary>",
			"hello world",
			"</details>",
		},
		details(
			[]string{
				"hello world",
			},
			"summary",
		),
	)
}

func TestTopicDetailsProvidesLink(t *testing.T) {
	assert.Equal(t,
		[]string{
			"<details>",
			"<summary>",
			"Broader related changes",
			"</summary>",
			"This change is part of a broader topic that can be in multiple repositories.",
			"<br/>",
			`Topic: <a href="https://search.example.com/topic" searchSha="89889b28e9672bff47fa4286f4aff4a80e09eade">some topic</a>`,
			"</details>",
		},
		topicDetails(&linkedTopicIssuesFunc{
			linkedTopicIssuesFunc: func(topicSearchString string) string {
				assert.Equal(t, "89889b28e9672bff47fa4286f4aff4a80e09eade", topicSearchString)
				return "https://search.example.com/topic"
			},
		}, "some topic"),
	)
}

type linkedTopicIssuesFunc struct {
	api.PullRequester
	linkedTopicIssuesFunc func(topicSearchString string) string
}

func (l linkedTopicIssuesFunc) LinkedTopicIssues(topicSearchString string) string {
	return l.linkedTopicIssuesFunc(topicSearchString)
}
