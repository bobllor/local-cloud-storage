package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type ServerConfiguration struct {
	Address string
}

func Read(path string) (*ServerConfiguration, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ServerConfiguration{}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
