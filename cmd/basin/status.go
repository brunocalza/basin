package main

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/brunocalza/basin/pkg/basin"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

var statusCmd = &cli.Command{
	Name:  "status",
	Usage: "get upload status",
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

		id, err := cid.Decode(cCtx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("decode cid: %s", err)
		}

		status, err := b.Status(id)
		if err != nil {
			return fmt.Errorf("list: %s", err)
		}

		s, _ := json.MarshalIndent(status, "", "\t")
		fmt.Printf("%s", string(s))

		return nil
	},
}
