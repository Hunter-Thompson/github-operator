package gh

import (
	"context"
	"os"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

func Login(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}
