package config

import (
	"os"
	"testing"

	"github.com/bobllor/assert"
	"github.com/goccy/go-yaml"
)

func TestReadConfig(t *testing.T) {
	s := newTestServerConfig()
	root := t.TempDir()

	buf, err := yaml.Marshal(s)
	assert.Nil(t, err)

	yamlPath, err := writeYaml(root, buf)
	assert.Nil(t, err)

	_, err = NewServerConfig(yamlPath)
	assert.Nil(t, err)
}

func writeYaml(root string, data []byte) (string, error) {
	name := "config.yml"
	yamlPath := root + "/" + name

	err := os.WriteFile(yamlPath, data, 0o644)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(yamlPath)
	if err != nil {
		return "", err
	}

	return yamlPath, nil
}

func newTestServerConfig() *ServerConfig {
	dbi := DatabaseInfo{
		Name:        "DBName",
		Address:     "127.0.0.1:3306",
		NetProtocol: "tcp",
		FileUser: DatabaseUser{
			User: "Username",
		},
		AccountUser: DatabaseUser{
			User: "UserUsername",
		},
	}

	s := &ServerConfig{
		Database: dbi,
		EnvFiles: []string{"test.env"},
	}

	return s
}
