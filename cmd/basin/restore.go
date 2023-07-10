package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/brunocalza/basin/pkg/basin"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

var restoreCmd = &cli.Command{
	Name:  "restore",
	Usage: "restore allows you to pick a snapshot to recover from",
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

		b, err := basin.NewBasin(cfg.DBS.Replica.AccessKey)
		if err != nil {
			return fmt.Errorf("creating basin impl: %s", err)
		}

		items, err := b.List()
		if err != nil {
			return fmt.Errorf("list: %s", err)
		}

		options := make([]string, len(items))
		for i, item := range items {
			options[i] = fmt.Sprintf("%s | %s", item.Created.Format(time.Stamp), item.Cid)
		}

		pick := ""
		prompt := &survey.Select{
			Message: "Choose a snapshot to recover from:",
			Options: options,
		}
		_ = survey.AskOne(prompt, &pick)

		parts := strings.Split(pick, " | ")
		id, err := cid.Decode(parts[1])
		if err != nil {
			return fmt.Errorf("decode cid: %s", err)
		}

		url, err := b.DownloadURL(id)
		if err != nil {
			return fmt.Errorf("get: %s", err)
		}

		dbstr := cCtx.Args().Get(0)
		dbFile, err := os.Create(dbstr)
		if err != nil {
			return fmt.Errorf("open file: %s", err)
		}
		defer dbFile.Close()

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("http get: %s", err)
		}
		defer resp.Body.Close()

		if _, err := io.Copy(dbFile, resp.Body); err != nil {
			return fmt.Errorf("copy file: %s", err)
		}

		return nil
	},
}
