package main

import (
	"context"
	"flag"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	comment := flag.String("comment", "", "A PR Comment")
	ghtoken := flag.String("ghtoken", "", "A github token")
	ghrepo := flag.String("ghrepo", "", "A github repository")
	prnumber := flag.String("prnumber", "", "A PR number")

	_ = ghrepo
	_ = prnumber
	_ = comment

	flag.Parse()

	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *ghtoken},
	)
	httpClient := oauth2.NewClient(context.Background(), token)

	client := githubv4.NewClient(httpClient)

	gh := Github{client: client}
	if err := gh.updatePrComment(context.Background(), 56, "einride", "prospective-shipper-service"); err != nil {
		panic(err)
	}

}
