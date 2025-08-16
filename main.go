package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	log *slog.Logger
)

func main() {
	var (
		logDebug bool
		dbFile   string
		err      error
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flags := flag.NewFlagSet("it488", flag.ContinueOnError)
	flags.StringVar(&dbFile, "db-file", "it488_team1.db", "Local path to sqlite database file")
	flags.BoolVar(&logDebug, "debug", false, "Enable debug logging")

	if err = flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			flags.PrintDefaults()
			os.Exit(0)
		}
		fmt.Print(err.Error())
		os.Exit(1)
	}

	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if logDebug {
		logOpts.Level = slog.LevelDebug
	}
	log = slog.New(slog.NewTextHandler(os.Stdout, logOpts))

	// spin up debug server
	go func() {
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	db, err := openDB(dbFile, logDebug)
	if err != nil {
		log.Error("Error opening database", "err", err)
		os.Exit(1)
	}

	fyneApp := app.New()
	logAppLifecycle(fyneApp)

	taskApp := newTaskApp(fyneApp, db)

	mainWindow := fyneApp.NewWindow("TODO Today")
	mainWindow.Resize(fyne.Size{Height: 700, Width: 300})
	mainWindow.SetContent(taskApp.Container())

	taskApp.RenderHomeView()

	// if context is cancelled, close app.
	go func() {
		<-ctx.Done()
		fyneApp.Quit()
	}()

	fyne.Do(mainWindow.ShowAndRun)
}
