package storer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

type GithubStorer struct {
	client *github.Client
	labels []string
	owner  string
	repo   string
}

// TODO: properly handle ratelimits

func (s *GithubStorer) Init(ctx context.Context) error {
	ghOwnerAndRepo := os.Getenv("GITHUB_REPOSITORY")

	token := os.Getenv("GITHUB_TOKEN")
	if len(ghOwnerAndRepo) == 0 || len(token) == 0 {
		return errors.New("env variables GITHUB_TOKEN and GITHUB_REPOSITORY must be set")
	}

	ghOwnerAndRepoParts := strings.Split(ghOwnerAndRepo, "/")
	if len(ghOwnerAndRepoParts) != 2 {
		return fmt.Errorf("GITHUB_REPOSITORY must be in the format owner/repo, got %s", ghOwnerAndRepo)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	rateLimit, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rate limit: %w", err)
	}
	if rateLimit.Core.Remaining < 1 {
		return fmt.Errorf("rate limit exceeded, remaining: %d", rateLimit.Core.Remaining)
	}

	*s = GithubStorer{
		client: client,
		labels: []string{"todo-syncer", "auto-generated"},
		owner:  ghOwnerAndRepoParts[0],
		repo:   ghOwnerAndRepoParts[1],
	}

	return nil
}

func (s GithubStorer) FetchCurrentOpenIssues(ctx context.Context) (map[string]Issue, error) {
	issues, _, err := s.client.Issues.ListByRepo(ctx, s.owner, s.repo, &github.IssueListByRepoOptions{
		State:  "open",
		Labels: s.labels,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues : %w", err)
	}

	result := make(map[string]Issue)
	for _, issue := range issues {
		body := issue.GetBody()

		firstLine := strings.Index(body, "\n")
		wannaBeID := strings.TrimSpace(body[:firstLine])
		if len(wannaBeID) == 0 {
			continue
		}
		if !strings.Contains(wannaBeID, ":") {
			slog.Warn("skipping issue %s, invalid format: %s", issue.GetHTMLURL(), wannaBeID)
			continue
		}

		result[wannaBeID] = Issue{
			Title:     issue.GetTitle(),
			Body:      body,
			ID:        strconv.Itoa(issue.GetNumber()),
			WannaBeID: wannaBeID,
		}
	}

	return result, nil
}

func (s GithubStorer) UpdateIssue(ctx context.Context, id string, issue Issue) error {
	issueRequest := &github.IssueRequest{
		Title:  &issue.Title,
		Body:   &issue.Body,
		Labels: &s.labels,
	}

	ghID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("failed to convert issue id %s : %w", id, err)
	}

	if ctx.Value(DryRunKey) != nil && ctx.Value(DryRunKey).(bool) {
		slog.Info("dry-run, not updating issue", "title", issue.Title, "id", id, "body", issue.Body)
		return nil
	}

	_, res, err := s.client.Issues.Edit(ctx, s.owner, s.repo, ghID, issueRequest)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to create issue %s : %w", issue.Title, err)
	}
	return nil
}

func (s GithubStorer) CreateIssue(ctx context.Context, issue Issue) error {
	issueRequest := &github.IssueRequest{
		Title:  &issue.Title,
		Body:   &issue.Body,
		Labels: &s.labels,
	}

	if ctx.Value(DryRunKey) != nil && ctx.Value(DryRunKey).(bool) {
		slog.Info("dry-run, not creating issue", "title", issue.Title, "body", issue.Body)
		return nil
	}

	_, res, err := s.client.Issues.Create(ctx, s.owner, s.repo, issueRequest)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to create issue %s : %w", issue.Title, err)
	}
	return nil
}

func (s GithubStorer) CloseIssue(ctx context.Context, id string) error {
	issueRequest := &github.IssueRequest{
		State: github.Ptr("closed"),
	}
	ghID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("failed to convert issue id %s : %w", id, err)
	}
	if ctx.Value(DryRunKey) != nil && ctx.Value(DryRunKey).(bool) {
		slog.Info("dry-run, not closing issue", "id", ghID)
		return nil
	}
	_, res, err := s.client.Issues.Edit(ctx, s.owner, s.repo, ghID, issueRequest)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to close issue %s : %w", id, err)
	}
	return nil
}
