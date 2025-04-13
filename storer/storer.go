package storer

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexdor/todo-syncer/parser"
)

type Issue struct {
	Title     string
	Body      string
	ID        string
	WannaBeID string
}

type Storer interface {
	FetchCurrentOpenIssues(ctx context.Context) (map[string]Issue, error)
	UpdateIssue(ctx context.Context, id string, issue Issue) error
	CreateIssue(ctx context.Context, issue Issue) error
	CloseIssue(ctx context.Context, id string) error
	Init(ctx context.Context) error
}

func generateAWannabeIDForComment(comment parser.Comment) string {
	return fmt.Sprintf("%s:%v", comment.FilePath, comment.LineNumber)
}

func IsValidWannabeID(id string) bool {
	parts := strings.Split(id, ":")

	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || !filepath.IsLocal(parts[0]) {
		return false
	}

	_, err := strconv.Atoi(parts[1])
	return err == nil
}

var AvailableStorer = map[string]Storer{
	"github": &GithubStorer{},
}

func UpdateIssues(ctx context.Context, s Storer, currentIssues map[string]Issue, comments []parser.Comment) error {
	for _, comment := range comments {
		title := comment.Text
		if len(title) > 120 {
			title = title[:120]
		}

		wannabeID := generateAWannabeIDForComment(comment)
		body := fmt.Sprintf("%s\n%s", wannabeID, comment.Text)
		if issue, ok := currentIssues[wannabeID]; ok {
			delete(currentIssues, wannabeID)
			if issue.Title == title && issue.Body == body {
				continue
			}
			err := s.UpdateIssue(ctx, issue.ID, issue)
			if err != nil {
				return fmt.Errorf("failed to update issue %s : %w", issue.ID, err)
			}
			continue
		}
		issue := Issue{
			Title:     title,
			WannaBeID: generateAWannabeIDForComment(comment),
			Body:      body,
		}
		err := s.CreateIssue(ctx, issue)
		if err != nil {
			return fmt.Errorf("failed to create issue %s %v:%v : %w",
				comment.FilePath, comment.LineNumber, comment.LineNumberEnd, err)
		}
	}
	for id := range currentIssues {
		err := s.CloseIssue(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete issue %s : %w", id, err)
		}
	}
	return nil
}
