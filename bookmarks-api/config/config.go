package config

import (
	"fmt"
	"os"

	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (*models.Config, error) {
	config := &models.Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)

	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}

	return nil
}
