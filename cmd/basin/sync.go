package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/urfave/cli/v2"
)

var syncCmd = &cli.Command{
	Name:  "sync",
	Usage: "periodically takes snapshots of your database and store it remotely",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Usage: "The directory where config is located (default: $HOME)",
		},
	},
	Action: func(cCtx *cli.Context) error {
		dir, err := defaultConfigLocation(cCtx)
		if err != nil {
			return fmt.Errorf("default config location: %s", err)
		}

		cfg, err := setupConfig(path.Join(dir, "config.yaml"))
		if err != nil {
			return fmt.Errorf("setup config: %s", err)
		}

		app, err := NewApp(cfg)
		if err != nil {
			log.Fatal(err)
		}

		app.Run()

		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done

		return nil
	},
}
