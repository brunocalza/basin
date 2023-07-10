package main

import (
	"fmt"
	"time"

	"github.com/brunocalza/basin/pkg/backup"
	"github.com/brunocalza/basin/pkg/basin"
)

type App struct {
	scheduler *backup.Scheduler
	uploader  *basin.Uploader
}

func NewApp(cfg *Config) (*App, error) {
	ch := make(chan backup.BackupResult)

	freq, err := time.ParseDuration(cfg.DBS.Replica.Frequency)
	if err != nil {
		return nil, fmt.Errorf("parse duration: %s", err)
	}

	backupScheduler, err := backup.NewScheduler(
		int(freq.Minutes()),
		ch,
		backup.BackuperOptions{
			SourcePath: cfg.DBS.Path,
			BackupDir:  "backups",
			Opts: []backup.Option{
				backup.WithVacuum(true),
			},
		})
	if err != nil {
		return nil, fmt.Errorf("creating backup scheduler: %s", err)
	}

	b, err := basin.NewBasin(cfg.DBS.Replica.AccessKey)
	if err != nil {
		return nil, fmt.Errorf("creating basin impl: %s", err)
	}

	return &App{
		scheduler: backupScheduler,
		uploader:  basin.NewUploader(ch, b),
	}, nil
}

func (app *App) Run() {
	go func() { app.scheduler.Run() }() // don't forget to close
	go func() { app.uploader.Start() }()
}
