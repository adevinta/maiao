package gh

import (
	"context"
	"fmt"
)

func ExampleNewGithubClient() {
	client, err := NewClient("github.com")
	if err != nil {
		Logger.Errorf("failed to create github client: %s", err.Error())
	}
	repo, _, err := client.Repositories.Get(context.Background(), "adevinta", "maiao")
	if err != nil {
		Logger.Errorf("failed to get github repository: %s", err.Error())
	}
	fmt.Println(repo.GetName())
	// Output: maiao
}
