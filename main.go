package main

import (
	"context"
	"crypto/sha1"
	"flag"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	ghComment := flag.String("comment", "", "A PR Comment")
	ghToken := flag.String("ghtoken", "", "A github token")
	ghRepo := flag.String("ghrepo", "", "A github repository")
	prNumber := flag.Int("prnumber", 0, "A Pull Request number")
	signKey := flag.String("signkey", "", "A key used to create the signature")

	flag.Parse()

	signature := sha1.New()
	_, err := signature.Write([]byte(*signKey))
	if err != nil {
		panic(err)
	}

	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), token)

	client := githubv4.NewClient(httpClient)

	gh := Github{client: client}
	if err := gh.CommentPR(context.Background(), *prNumber, "einride", *ghRepo, *ghComment, fmt.Sprintf("%x", signature.Sum(nil))); err != nil {
		panic(err)
	}

}
