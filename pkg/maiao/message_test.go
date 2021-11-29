package maiao

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/adevinta/maiao/pkg/api"
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
			`Topic: <a href="https://github.company.example.com/search?q=type%3Apr+%22Topic%3A+some+topic%22&type=Issues">some topic</a>`,
			"</details>",
		},
		topicDetails(&api.GitHub{Host: "github.company.example.com"}, "some topic"),
	)
}
