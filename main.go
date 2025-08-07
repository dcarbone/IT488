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

	"fyne.io/fyne/v2/app"
)

var (
	logDebug bool
	log      *slog.Logger
)

func main() {
	var (
		err error
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flags := flag.NewFlagSet("it488", flag.ContinueOnError)
	flags.StringVar(&dbFile, "db-file", "it488_group1.db", "Local path to sqlite database file")
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

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	db, err := openDB()
	if err != nil {
		log.Error("Error opening database", "err", err)
		os.Exit(1)
	}

	a := app.New()
	logLifecycle(a)
	w := a.NewWindow("Daniel Carbone IT481")
	w.SetOnClosed(func() { cancel() })
}
