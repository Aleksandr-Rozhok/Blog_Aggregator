package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"
const projectPath = "/Desktop/Programming/Projects/Go/Blog_Aggregator/"

type Config struct {
	DBURL           string "json:db_url"
	CurrentUserName string "json:current_user_name"
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return nil, err
	}

	return &config, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "Wrong file's path", err
	}
	return homePath + projectPath + configFileName, nil
}

func write(cfg *Config) error {
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(c)
}
