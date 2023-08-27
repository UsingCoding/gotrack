package tracker

import (
	"context"
	stderrors "errors"
	"time"
)

var (
	ErrIssueNotFound = stderrors.New("issue not found")
)

type TrackParams struct {
	IssueID  string
	Text     string
	Duration time.Duration
	Date     time.Time
}

type Client interface {
	IssueExists(ctx context.Context, issueID string) (bool, error)
	Track(ctx context.Context, params TrackParams) error
}
