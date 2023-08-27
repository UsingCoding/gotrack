package dataprovider

import (
	"context"
	"time"
)

type Record struct {
	IssueID string
	Comment string
	Start   time.Time
	End     time.Time
	Tracked bool
}

type Provider interface {
	Read(ctx context.Context) ([]Record, error)
}
