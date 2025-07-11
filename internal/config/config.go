package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigItem struct {
	Name    string
	Command string
}

var printers = map[int]ConfigItem{
	1: {"Salita", "lpr -P hp-335"},
	2: {"Toqui", "lpr"},
}

var modes = map[int]ConfigItem{
	1: {"Simple", ""},
	2: {"Doble", "duplex"},
}

var borders = map[int]ConfigItem{
	1: {"Corto", "duplex -l"},
	2: {"Largo", "duplex"},
}

// Todo: support -dFirstPage= y -dLastPage= from postscript (or psselect -p5-10)
// Todo: Consultar papel?
type Config struct {
	Theme          string `json:"theme"`
	Account        string `json:"account"`
	Printer        string `json:"printer"`
	Mode           string `json:"mode"`
	Border         string `json:"border"`
	SetupCompleted bool   `json:"setupCompleted"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dccprint_config.json"), nil
}

func Load() Config {
	defaultConfig := Config{Theme: "Default", Account: "", Printer: "Toqui", Mode: "Doble", Border: "Largo"}
	path, err := configPath()
	if err != nil {
		return defaultConfig
	}

	file, err := os.Open(path)
	if err != nil {
		return defaultConfig
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return defaultConfig
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

func updateConfig(updater func(cfg *Config)) error {
	cfg := Load()
	updater(&cfg)
	return save(cfg)
}

func SaveConfig(cfg Config) error {
	return save(cfg)
}

func SaveTheme(theme string) error {
	return updateConfig(func(cfg *Config) { cfg.Theme = theme })
}

func SaveAccount(account string) error {
	return updateConfig(func(cfg *Config) { cfg.Account = account })
}

func SavePrinter(printer string) error {
	return updateConfig(func(cfg *Config) { cfg.Printer = printer })
}

func SaveMode(mode string) error {
	return updateConfig(func(cfg *Config) { cfg.Mode = mode })
}

func SaveBorder(border string) error {
	return updateConfig(func(cfg *Config) { cfg.Border = border })
}
