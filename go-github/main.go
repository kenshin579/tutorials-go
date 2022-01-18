package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "11"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	fmt.Printf("err: %+v\n", err)
	fmt.Printf("repos: %+v\n", repos)
}
