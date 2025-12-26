package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DatabaseURL string `json:"db_url"`
	User        string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homePath + "/" + configFileName, nil
}

func write(config *Config) error {

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Read() (*Config, error) {

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SetUser(user string) error {
	c.User = user
	return write(c)
}
