package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Printer struct {
	Name    string
	Command string
}

var printers = map[int]Printer{
	1: {"Salita", "lpr -P hp-335"},
	2: {"Toqui", "lpr"},
}

type Config struct {
	Theme   string `json:"theme"`
	Account string `json:"account"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dccprint_config.json"), nil
}

func Load() Config {
	path, err := configPath()
	if err != nil {
		return Config{Theme: "Default", Account: ""}
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{Theme: "Default", Account: ""}
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{Theme: "Default", Account: ""}
	}

	if cfg.Theme == "" {
		cfg.Theme = "Default"
	}

	return cfg
}

func save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func SaveTheme(theme string) error {
	cfg := Load()
	cfg.Theme = theme
	return save(cfg)
}

func SaveAccount(account string) error {
	cfg := Load()
	cfg.Account = account
	return save(cfg)
}
