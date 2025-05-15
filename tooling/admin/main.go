package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/tooling/admin/commands"
	"github.com/ardanlabs/conf/v3"
)

var build = "develop"

type config struct {
	conf.Version
	Args conf.Args
	Auth struct {
		KeysFolder string `conf:"default:zarf/keys/"`
	}
}

func main() {
	log := logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	if err := run(log); err != nil {
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "development build",
		},
	}

	const prefix = "ADMIN"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Info(context.Background(), "startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands executes the command provided on the command line
func processCommands(args conf.Args, log *logger.Logger, cfg config) error {
	switch args.Num(0) {
	case "genkey":
		if err := commands.GenKey(cfg.Auth.KeysFolder); err != nil {
			return fmt.Errorf("generating key: %w", err)
		}
	}
	return nil
}
