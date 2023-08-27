package app

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gotrack/pkg/app/dataprovider"
	appreporter "gotrack/pkg/app/reporter"
	"gotrack/pkg/app/tracker"
)

type Service interface {
	Track(ctx context.Context, issueID, comment string, spentTime time.Duration, date time.Time) error
	TrackTimeFromTable(ctx context.Context, dryRun bool) error
}

func NewService(
	tracker tracker.Client,
	provider dataprovider.Provider,
	reporter appreporter.Reporter,
) Service {
	return &service{
		tracker:      tracker,
		dataProvider: provider,
		reporter:     reporter,
	}
}

type service struct {
	tracker      tracker.Client
	dataProvider dataprovider.Provider
	reporter     appreporter.Reporter
}

func (s *service) Track(ctx context.Context, issueID, comment string, spentTime time.Duration, date time.Time) error {
	err := s.tracker.Track(ctx, tracker.TrackParams{
		IssueID:  issueID,
		Text:     comment,
		Duration: spentTime,
		Date:     date,
	})
	if err != nil {
		return err
	}

	s.reporter.Report(fmt.Sprintf("Successfully tracked time for %s : %s", issueID, spentTime.String()))
	return err
}

func (s *service) TrackTimeFromTable(ctx context.Context, dryRun bool) error {
	records, err := s.dataProvider.Read(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(records) == 0 {
		s.reporter.Report("Empty records, there is nothing to track")
		return nil
	}

	var totalTime time.Duration
	for _, record := range records {
		if record.Tracked {
			s.reporter.Report(fmt.Sprintf("TimeTrack record for Issue %s marked as 'tracked', skip", record.IssueID))
			continue
		}

		spentTime := record.End.Sub(record.Start)
		totalTime += spentTime

		var issueNotFound bool

		if !dryRun {
			err2 := s.tracker.Track(ctx, tracker.TrackParams{
				IssueID:  record.IssueID,
				Text:     record.Comment,
				Duration: spentTime,
				Date:     time.Now(),
			})

			switch errors.Cause(err2) {
			case tracker.ErrIssueNotFound:
				issueNotFound = true
			default:
				return errors.Wrapf(err2, "failed track time for issue %s", record.IssueID)
			case nil:
				// leave switch on nil
			}
		}

		var postMsg string
		switch {
		case dryRun:
			postMsg = ", dry-run"
		case issueNotFound:
			postMsg = ", skipped - issue not found"
		}

		s.reporter.Report(fmt.Sprintf(
			"Tracked time for %s, spent time %s%s",
			record.IssueID,
			spentTime,
			postMsg,
		))
	}

	s.reporter.Report(fmt.Sprintf("Total tracked time: %s", totalTime))

	return nil
}
