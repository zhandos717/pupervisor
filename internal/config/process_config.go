package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ProcessConfig struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args,omitempty"`
	Directory   string            `yaml:"directory,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	AutoStart   bool              `yaml:"autostart"`
	AutoRestart bool              `yaml:"autorestart"`
	StartSecs   int               `yaml:"startsecs,omitempty"`
	StopSignal  string            `yaml:"stopsignal,omitempty"`
	StopTimeout int               `yaml:"stoptimeout,omitempty"`
	Stdout      string            `yaml:"stdout,omitempty"`
	Stderr      string            `yaml:"stderr,omitempty"`
}

type SupervisorConfig struct {
	Processes []ProcessConfig `yaml:"processes"`
}

func LoadProcessConfig(path string) (*SupervisorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg SupervisorConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	for i := range cfg.Processes {
		if cfg.Processes[i].StopSignal == "" {
			cfg.Processes[i].StopSignal = "SIGTERM"
		}
		if cfg.Processes[i].StopTimeout == 0 {
			cfg.Processes[i].StopTimeout = 10
		}
		if cfg.Processes[i].StartSecs == 0 {
			cfg.Processes[i].StartSecs = 1
		}
	}

	return &cfg, nil
}
