package config

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const settingsPath = "/app/config"

var settingsFiles = map[string]string{
	"development": path.Join(settingsPath, "appConfig.dev.yml"),
	"production":  path.Join(settingsPath, "appConfig.yml"),
}

// DBSettings contains the settings used for ths database connection
type DBSettings struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Settings struct {
	DB DBSettings
}

var environment string = os.Getenv("Environment")

// LoadSettings loads the settings from the yml file
func LoadSettings() (*Settings, error) {
	config, err := ioutil.ReadFile(settingsFiles[environment])
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	err = yaml.Unmarshal(config, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
