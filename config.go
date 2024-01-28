package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/GioPan04/bereal"
)

type Config struct {
	RefreshAt *time.Time
	Session   *bereal.BeRealSession
}

func LoadConfig(file string) (*Config, error) {
	fi, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(fi, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save(file string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, data, 0660)
	if err != nil {
		return err
	}

	return nil
}
