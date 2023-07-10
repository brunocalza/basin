package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var newCmd = &cli.Command{
	Name:  "new",
	Usage: "initializes basin",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Usage: "The directory where config will be stored (default: $HOME)",
		},
	},
	Action: func(cCtx *cli.Context) error {
		dir, err := defaultConfigLocation(cCtx)
		if err != nil {
			return fmt.Errorf("default config location: %s", err)
		}

		cfg := Config{}
		f, err := os.Create(path.Join(dir, "config.yaml"))
		if err != nil {
			return fmt.Errorf("os create: %s", err)
		}

		answers := struct {
			Path      string
			AccessKey string
			Frequency string
		}{}

		// the questions to ask
		qs := []*survey.Question{
			{
				Name: "path",
				Prompt: &survey.Input{
					Message: "Enter your database path: ",
				},
				Validate: survey.Required,
			},
			{
				Name: "accessKey",
				Prompt: &survey.Input{
					Message: "Enter you Basin secret key:",
				},
				Validate: survey.Required,
			},
			{
				Name: "frequency",
				Prompt: &survey.Input{
					Message: "Snapshot frequency",
					Default: "24h",
				},
				Validate: survey.Required,
			},
		}

		if err := survey.Ask(qs, &answers); err != nil {
			return fmt.Errorf("survey ask: %s", err)
		}

		connectionID, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("new uuid: %s", err)
		}

		cfg.DBS.Path = answers.Path
		cfg.DBS.Replica.AccessKey = answers.AccessKey
		cfg.DBS.Replica.Frequency = answers.Frequency
		cfg.DBS.Replica.URL = fmt.Sprintf("https://basin.tableland.xyz/%s", strings.ReplaceAll(connectionID.String(), "-", ""))

		if err := yaml.NewEncoder(f).Encode(cfg); err != nil {
			return fmt.Errorf("encode: %s", err)
		}

		fmt.Printf("\n\n\033[32mSuccess!\033[0m\n\n")

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshal: %s", err)
		}

		fmt.Println(string(bytes))
		fmt.Println()
		fmt.Printf("\033[32mwritten to %s\033[0m\n\n", path.Join(dir, "config.yaml"))

		return nil
	},
}
