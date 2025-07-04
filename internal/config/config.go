package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Theme string `json:"theme"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dccprint_config.json")
}

func Load() Config {
	path := configPath()
	file, err := os.Open(path)
	if err != nil {
		return Config{Theme: "Default"}
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{Theme: "Default"}
	}
	return cfg
}

func Save(theme string) {
	path := configPath()
	cfg := Config{Theme: theme}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	_ = encoder.Encode(cfg)
}
