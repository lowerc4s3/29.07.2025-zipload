package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port         int      `json:"port"`
	AllowedMIMEs []string `json:"allowed_mimes"`
	MaxSources   int      `json:"max_sources"`
	MaxTasks     int      `json:"max_tasks"`
	MaxTaskFiles int      `json:"max_task_files"`
}

func ReadConfig() (Config, error) {
	cfgJSON, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
