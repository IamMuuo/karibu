package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure for the application
type Config struct {
	// Application configuration
	AppConfig struct {
		Port   int    `yaml:"port"`
		Addres string `yaml:"address"`
	} `yaml:"server"`
}

// LoadDefaultConfigs loads the default configuration for the application.
// It returns a Config instance upon success and an error in the event of failure
func LoadDefaultConfigs() (*Config, error) {
	var cfg Config

	if err := processYAMLConfigs(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// The processYAMLConfigs loads a configuration from a YAML file
// It requires a Config pointer to load data onto it and returns an error
// incase loading failed
func processYAMLConfigs(cfg *Config) error {

	f, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	// Parsing the yaml file
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}
