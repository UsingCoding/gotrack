package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	appID = "gotrack"
)

func main() {
	ctx := context.Background()

	err := runApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp(
	ctx context.Context,
) (err error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	ctx = listenOSKillSignals(ctx)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "fetch home dir")
	}

	app := cli.App{
		Name:  appID,
		Usage: "Time tracking helper",
		Commands: []*cli.Command{
			track(homeDir),
		},
	}
	err = app.RunContext(ctx, os.Args)
	return err
}

// Subscribes for os kill signals and cancel context on notification
func listenOSKillSignals(ctx context.Context) context.Context {
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(ctx)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ch:
			cancelFunc()
		case <-ctx.Done():
			signal.Reset()
			return
		}
	}()

	return ctx
}
