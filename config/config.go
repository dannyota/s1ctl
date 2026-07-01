// Package config locates, loads, and validates an s1ctl instance configuration.
//
// Per-value resolution, highest priority first:
//
//  1. S1_* environment variables
//  2. config file at an explicit path (--config flag)
//  3. ~/.s1ctl/config.yaml (default)
//  4. ./config/config.yaml (local fallback)
//
// A file value is overlaid by the matching env var when set.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Instance is one SentinelOne console configuration.
type Instance struct {
	ConsoleURL string `yaml:"console_url"`
	Token      string `yaml:"token"`
	SDLURL     string `yaml:"sdl_url"`

	source string
}

// Source returns the config file this instance was loaded from, or "env" when
// loaded from environment variables only.
func (i *Instance) Source() string {
	if i.source == "" {
		return "env"
	}
	return i.source
}

// DefaultPath returns ~/.s1ctl/config.yaml.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".s1ctl", "config.yaml")
}

// Load resolves an Instance from the config chain (env overlays file).
// If explicit is non-empty, only that file path is tried.
func Load(explicit string) (*Instance, error) {
	inst := &Instance{}

	for _, p := range filePaths(explicit) {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if err := yaml.Unmarshal(data, inst); err != nil {
			return nil, fmt.Errorf("config: %s: %w", p, err)
		}
		inst.source = p
		break
	}

	if v := os.Getenv("S1_CONSOLE_URL"); v != "" {
		inst.ConsoleURL = v
	}
	if v := os.Getenv("S1_TOKEN"); v != "" {
		inst.Token = v
	}
	if v := os.Getenv("S1_SDL_URL"); v != "" {
		inst.SDLURL = v
	}

	inst.ConsoleURL = strings.TrimRight(inst.ConsoleURL, "/")
	inst.SDLURL = strings.TrimRight(inst.SDLURL, "/")
	return inst, nil
}

// Validate checks that the minimum required fields are set.
func (i *Instance) Validate() error {
	if i.ConsoleURL == "" {
		return fmt.Errorf("config: console URL is required (set S1_CONSOLE_URL or console_url in config)")
	}
	if i.Token == "" {
		return fmt.Errorf("config: API token is required (set S1_TOKEN or token in config)")
	}
	return nil
}

// Save writes the instance config to path, creating parent dirs as needed.
// The file is written with 0600 permissions (token may be present).
func Save(path string, inst *Instance) error {
	data, err := yaml.Marshal(inst)
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("config: mkdir: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// ResolvedSource returns the path of the config file that would be loaded,
// or "" if none exists.
func ResolvedSource(explicit string) string {
	for _, p := range filePaths(explicit) {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// ReadForEdit loads the config at path for interactive editing, returning an
// empty Instance if the file doesn't exist or is unreadable.
func ReadForEdit(path string) *Instance {
	inst, err := Load(path)
	if err != nil || inst == nil {
		return &Instance{}
	}
	return inst
}

func filePaths(explicit string) []string {
	if explicit != "" {
		return []string{explicit}
	}
	var paths []string
	if p := DefaultPath(); p != "" {
		paths = append(paths, p)
	}
	paths = append(paths, filepath.Join("config", "config.yaml"))
	return paths
}
