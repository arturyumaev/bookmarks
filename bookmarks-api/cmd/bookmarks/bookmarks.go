package main

import (
	"flag"

	"github.com/arturyumaev/bookmarks/bookmarks-api/config"
	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/app"
	"github.com/sirupsen/logrus"
)

type Flags struct {
	ConfigPath string
}

func parseFlags() (*Flags, error) {
	flags := &Flags{}

	var configPath string
	flag.StringVar(&configPath, "config", "./config/config.yaml", "path to config file")
	flag.Parse()

	if err := config.ValidateConfigPath(configPath); err != nil {
		return nil, err
	}

	flags.ConfigPath = configPath

	return flags, nil
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		logrus.Fatal(err)
	}

	config, err := config.ReadConfig(flags.ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	app := app.NewApplication(config)
	defer app.BoltDB.Close()

	if err = app.Run(); err != nil {
		logrus.Fatal(err)
	}
}
