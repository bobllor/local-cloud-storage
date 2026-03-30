package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Database         DatabaseInfo `yaml:"database"`    // DatabaseInfo holds the info for establishing a database connection.
	EnvironmentFiles []string     `yaml:"environment"` // EnvironmentFiles holds the .env file paths to be loaded into the program.
}

type DatabaseInfo struct {
	Name            string       `yaml:"name"`             // Name is the database name to connect to.
	Address         string       `yaml:"address"`          // Address is the address of the database. This includes the port.
	NetProtocol     string       `yaml:"network_protocol"` // NetProtocol is the network protocol used for the connection.
	FileUser        DatabaseUser `yaml:"file_user"`        // FileUser represents the user that handles the File table.
	UserAccountUser DatabaseUser `yaml:"useraccount_user"` // UserAccountUser represents the user that handles the UserAccount table.
}

type DatabaseUser struct {
	// User is the username account used for the connection.
	User string `yaml:"username"`
	// Password is the user password used for the connection. This can be loaded
	// into the env file but it must use the correct value.
	Password string `yaml:"user_password"`
}

// NewServerConfig creates a new ServerConfig read from a path.
func NewServerConfig(path string) (*ServerConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ServerConfig{}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadEnv loads the .env files into the program.
func (s *ServerConfig) LoadEnv() error {
	err := godotenv.Load(s.EnvironmentFiles...)
	if err != nil {
		return fmt.Errorf("failed to load .env files: %v", err)
	}

	return nil
}

// LoadDbSecrets loads the environment variables used for the
// database account.
//
// If an environment variable exists for the related variable,
// then it will replace the given value inside the information
// only if the key is used for the variable.
func (s *ServerConfig) LoadDbSecrets() {
	file_key := "FILEUSER_PASSWORD"
	useraccount_key := "USERACCOUNTUSER_PASSWORD"

	file_pw := os.Getenv(file_key)
	useraccount_pw := os.Getenv(useraccount_key)

	if s.Database.FileUser.Password == "$"+file_key {
		s.Database.FileUser.Password = file_pw
	}

	if s.Database.UserAccountUser.Password == "$"+useraccount_key {
		s.Database.UserAccountUser.Password = useraccount_pw
	}
}
