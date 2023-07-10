package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/brunocalza/basin/pkg/basin"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var listCmd = &cli.Command{
	Name:  "list",
	Usage: "list lists available backups",
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

		data := make([][]string, len(items))
		for i, item := range items {
			data[i] = make([]string, 2)
			data[i][0] = item.Created.Format(time.Stamp)
			data[i][1] = item.Cid.String()
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Snapshots", "Links"})
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()

		return nil
	},
}
