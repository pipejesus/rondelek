package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Pads   []PadConfig `json:"pads"`
	Window Window      `json:"window"`
	Layout Layout      `json:"layout"`
}

type Window struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type Layout struct {
	Columns  int     `json:"columns"`
	Rows     int     `json:"rows"`
	PaddingX float32 `json:"padding_x"`
	PaddingY float32 `json:"padding_y"`
}

type PadConfig struct {
	Type        string      `json:"type"`
	ID          string      `json:"id"`
	Label       string      `json:"label"`
	Key         int32       `json:"key"`
	PadPosition PadPosition `json:"position"`
	PadSize     PadSize     `json:"size"`
}

type PadPosition struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

type PadSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewConfig() *Config {
	return &Config{
		Pads: []PadConfig{},
	}
}

// Load reads the config.json file from the program's directory (i.e., the
// directory where the executable resides) and unmarshals it into the receiver.
// If the file is not found in the executable's directory, it will try to
// load config.json from the current working directory as a fallback.
func (c *Config) Load() error {
	exePath, err := os.Executable()
	if err != nil {
		if err2 := c.loadFromCurrentDir(); err2 != nil {
			return fmt.Errorf("Error loading config file: cannot resolve the executable path (%v) and failed to load config from your current directory (%v)", err, err2)
		}
		return nil
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")
	if err := c.loadFromFile(configPath); err == nil {
		return nil
	}

	if err := c.loadFromCurrentDir(); err != nil {
		return fmt.Errorf("Error loading config file: tried %s and working directory: %w", configPath, err)
	}
	return nil
}

func (c *Config) loadFromCurrentDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error loading config: failed to get current working directory: %w", err)
	}
	path := filepath.Join(cwd, "config.json")
	return c.loadFromFile(path)
}

func (c *Config) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error loading config: failed to read config file at %s: %w", path, err)
	}
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("Error loading config: failed to parse config JSON at %s: %w", path, err)
	}
	return nil
}
