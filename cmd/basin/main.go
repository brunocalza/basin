package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	zlog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()

	cliApp := &cli.App{
		Name:  "basin",
		Usage: "basin creates snapshots your SQLite database periodically and store them in Filecoin",
		Commands: []*cli.Command{
			newCmd,
			syncCmd,
			listCmd,
			restoreCmd,
			statusCmd,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func defaultConfigLocation(cCtx *cli.Context) (string, error) {
	var dir string
	dir = cCtx.String("dir")
	if dir == "" {
		// the default directory is home
		var err error
		dir, err = homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("home dir: %s", err)
		}

		dir = path.Join(dir, ".basin")
	}

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0o755); err != nil {
			return "", fmt.Errorf("mkdir: %s", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("is not exist: %s", err)
	}

	return dir, nil
}
