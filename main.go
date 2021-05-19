package main

import (
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	comment := flag.String("comment", "", "A PR Comment")
	ghToken := flag.String("token", "", "A github token")
	repo := flag.String("repo", "", "A github repository. In the format einride/<repository>. If empty the local git repo will be used")
	pr := flag.Int("pr", 0, "A Pull Request number")
	signKey := flag.String("signkey", "", "A key used to create the signature")

	flag.Parse()

	if *repo == "" {
		path, err := retrieveLocalGitRepo()
		if err != nil {
			panic(fmt.Errorf("failed to parse local git repository: %v", err))
		}
		*repo = path
	}
	owner, name, err := parseGitPath(*repo)
	if err != nil {
		panic(err)
	}

	signature := sha1.New()
	_, err = signature.Write([]byte(*signKey))
	if err != nil {
		panic(err)
	}

	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), token)

	client := githubv4.NewClient(httpClient)

	gh := Github{client: client}
	if err := gh.CommentPR(context.Background(), *pr, owner, name, *comment, fmt.Sprintf("%x", signature.Sum(nil))); err != nil {
		panic(err)
	}
}

func parseGitPath(path string) (owner string, name string, _ error) {
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid repository format %s", path)
	}
	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}

func retrieveLocalGitRepo() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to open .git directory: %v", err)
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("failed to list remotes: %v", err)
	}
	if len(remotes) != 1 {
		return "", fmt.Errorf("found %d remotes expecting 1", len(remotes))
	}
	if err := remotes[0].Config().Validate(); err != nil {
		return "", fmt.Errorf("invalid remote config: %v", err)
	}
	endpoint, err := transport.NewEndpoint(remotes[0].Config().URLs[0])
	if err != nil {
		return "", fmt.Errorf("unable to parse remote endpoint: %v", err)
	}
	if endpoint.Host != "github.com" {
		return "", fmt.Errorf("not a github.com repository")
	}
	return endpoint.Path, nil
}
