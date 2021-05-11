package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
)

type Github struct {
	client *githubv4.Client
}

func (g *Github) CommentPR(ctx context.Context, prnumber int, owner, name, body, signature string) error {
	body = body + "\n <sub>Sign: " + signature + "</sub>"

	type Comment struct {
		ID   githubv4.ID
		Body string
	}

	var prComments struct {
		Repository struct {
			PullRequest struct {
				ID       githubv4.ID
				Comments struct {
					Nodes    []Comment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"`
			} `graphql:"pullRequest(number: $prnumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":          githubv4.String(owner),
		"name":           githubv4.String(name),
		"prnumber":       githubv4.Int(prnumber),
		"commentsCursor": (*githubv4.String)(nil),
	}

	var allComments []Comment
	var comment *Comment
LOOP:
	for {
		err := g.client.Query(ctx, &prComments, variables)
		if err != nil {
			return fmt.Errorf("Unable to query PR: %v", err)
		}
		for _, c := range prComments.Repository.PullRequest.Comments.Nodes {
			c := c
			if strings.Contains(c.Body, signature) {
				comment = &c
				break LOOP
			}
		}
		allComments = append(allComments, prComments.Repository.PullRequest.Comments.Nodes...)
		if !prComments.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(prComments.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	if comment == nil {
		if err := g.addComment(ctx, prComments.Repository.PullRequest.ID, body); err != nil {
			return fmt.Errorf("unable to generate new comment on PR: %v", err)
		}
		return nil
	}
	if err := g.updateComment(ctx, body, comment.ID); err != nil {
		return fmt.Errorf("unable to update comment on PR: %v", err)
	}

	return nil
}

func (g *Github) updateComment(ctx context.Context, body string, commentID githubv4.ID) error {
	fmt.Println("Updating comment")

	var updateIssueComment struct {
		UpdateIssueComment struct {
			IssueComment struct {
				ID githubv4.ID
			}
		} `graphql:"updateIssueComment(input: $input)"`
	}

	updateIssueCommentInput := githubv4.UpdateIssueCommentInput{
		ID:   commentID,
		Body: githubv4.String(body),
	}

	err := g.client.Mutate(ctx, &updateIssueComment, updateIssueCommentInput, nil)
	if err != nil {
		return fmt.Errorf("Unable to update comment in PR: %v", err)
	}
	return nil
}

func (g *Github) addComment(ctx context.Context, prID githubv4.ID, body string) error {
	fmt.Println("Creating new comment")

	addCommentInput := githubv4.AddCommentInput{
		SubjectID: prID,
		Body:      githubv4.String(body),
	}

	var addComment struct {
		AddComment struct {
			Subject struct {
				ID githubv4.ID
			}
		} `graphql:"addComment(input: $input)"`
	}

	err := g.client.Mutate(ctx, &addComment, addCommentInput, nil)
	if err != nil {
		return fmt.Errorf("Unable to comment in PR: %v", err)
	}
	return nil
}
