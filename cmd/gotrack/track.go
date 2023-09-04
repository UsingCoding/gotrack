package main

import (
	"log"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gotrack/pkg/app"
	"gotrack/pkg/infrastructure/config"
	"gotrack/pkg/infrastructure/dataprovider"
	"gotrack/pkg/infrastructure/reporter"
	"gotrack/pkg/infrastructure/youtrack"
)

func track(homeDir string) *cli.Command {

	return &cli.Command{
		Name:   "track",
		Action: executeTrack,
		Usage:  "Manually track time: <issueID> <comment> <spentTime> <date>",
		Subcommands: []*cli.Command{
			{
				Name:   "table",
				Usage:  "Track time from table",
				Action: executeTrackTable,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "timetrack",
						Aliases: []string{"tt"},
						Usage:   "Path to csv file with time track",
					},
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "Date of time track",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Run track in dry-run",
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to gotrack config",
				Value:   path.Join(homeDir, ".gotrack/config"),
			},
		},
	}
}

func executeTrack(ctx *cli.Context) error {
	configPath := ctx.String("config")
	if configPath == "" {
		return errors.New("empty path to config")
	}
	c, err := config.Parser{}.Parse(configPath)
	if err != nil {
		return err
	}

	argsLen := ctx.Args().Len()
	if argsLen < 3 || argsLen > 4 {
		return errors.Errorf("invalid count of args %d, expects 3 or 4", argsLen)
	}

	issueID := ctx.Args().Get(0)
	comment := ctx.Args().Get(1)
	spentTime, err := time.ParseDuration(ctx.Args().Get(2))
	if err != nil {
		return errors.Wrapf(err, "failed to parse spentTime %s", ctx.Args().Get(2))
	}

	date := time.Now()
	if dateStr := ctx.Args().Get(3); dateStr != "" {
		date, err = time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse date %s", dateStr)
		}
	}

	tracker := youtrack.NewClient(
		c.YouTrackHost,
		c.Token,
	)
	srv := app.NewService(
		tracker,
		nil, // pass nil as provider
		reporter.NewReporter(log.Default()),
	)

	return srv.Track(ctx.Context, issueID, comment, spentTime, date)
}

func executeTrackTable(ctx *cli.Context) error {
	timeTrackFilePath := ctx.String("timetrack")
	if timeTrackFilePath == "" {
		return errors.New("empty path to time track file")
	}

	dateStr := ctx.String("date")
	if dateStr == "" {
		return errors.New("empty date to time track")
	}

	date := time.Now()
	if dateStr != "" {
		var err error
		date, err = time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse date %s", dateStr)
		}
	}
	configPath := ctx.String("config")
	if configPath == "" {
		return errors.New("empty path to config")
	}
	c, err := config.Parser{}.Parse(configPath)
	if err != nil {
		return err
	}

	dryRun := ctx.Bool("dry-run")

	tracker := youtrack.NewClient(
		c.YouTrackHost,
		c.Token,
	)
	srv := app.NewService(
		tracker,
		dataprovider.NewCsvFileDataProvider(timeTrackFilePath),
		reporter.NewReporter(log.Default()),
	)

	return srv.TrackTimeFromTable(ctx.Context, date, dryRun)
}
