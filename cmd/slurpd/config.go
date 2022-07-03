package main

import "encoding/json"

type Network struct {
	Address string
	User    string
	Pass    string
}

type Config struct {
	Networks []Network
}

// Load creates a `Config` from a JSON string
func LoadConfig(j string) (*Config, error) {
	var c Config
	if err := json.Unmarshal([]byte(j), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
