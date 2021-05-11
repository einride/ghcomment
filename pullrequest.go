package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type Github struct {
	client *githubv4.Client
}

func (g *Github) updatePrComment(ctx context.Context, prnumber int, owner, name string) error {

	type comment struct {
		Id   githubv4.ID
		Body string
	}

	var prComments struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes    []comment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"`
			} `graphql:"pullRequest(number: $prnumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	//{
	//	repository(name: "", owner: "") {
	//	  pullRequest(number: 10) {
	//		comments {
	//		  nodes {
	//			body
	//			id
	//		  }
	//		}
	//	  }
	//	}
	//  }

	variables := map[string]interface{}{
		"owner":          githubv4.String(owner),
		"name":           githubv4.String(name),
		"prnumber":       githubv4.Int(prnumber),
		"commentsCursor": (*githubv4.String)(nil),
	}

	var allComments []comment
	for {
		err := g.client.Query(context.Background(), &prComments, variables)
		if err != nil {
			return fmt.Errorf("Unable to query PR: %v", err)
		}
		allComments = append(allComments, prComments.Repository.PullRequest.Comments.Nodes...)
		if !prComments.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(prComments.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	fmt.Printf("    PrComment: %+v \n", allComments)

	fmt.Println("Length: ", len(allComments))

	return nil

}
