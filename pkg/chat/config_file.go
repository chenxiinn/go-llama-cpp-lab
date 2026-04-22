package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const DefaultConfigPath = "config/local.json"

type runtimeConfigFile struct {
	ModelPath   *string `json:"model_path"`
	ContextSize *int    `json:"n_ctx"`
	GPULayers   *int    `json:"gpu_layers"`
	MaxTokens   *int    `json:"max_tokens"`
}

func ResolveConfigPath(args []string, defaultPath string) (string, bool) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-config" || arg == "--config":
			if i+1 < len(args) {
				return args[i+1], true
			}
		case strings.HasPrefix(arg, "-config="):
			return strings.TrimPrefix(arg, "-config="), true
		case strings.HasPrefix(arg, "--config="):
			return strings.TrimPrefix(arg, "--config="), true
		}
	}

	return defaultPath, false
}

func (c *RuntimeConfig) ApplyOptionalJSONFile(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat config file %q: %w", path, err)
	}
	return c.ApplyJSONFile(path)
}

func (c *RuntimeConfig) ApplyJSONFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", path, err)
	}

	var fileCfg runtimeConfigFile
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return fmt.Errorf("decode config file %q: %w", path, err)
	}

	if fileCfg.ModelPath != nil {
		c.ModelPath = *fileCfg.ModelPath
	}
	if fileCfg.ContextSize != nil {
		c.ContextSize = *fileCfg.ContextSize
	}
	if fileCfg.GPULayers != nil {
		c.GPULayers = *fileCfg.GPULayers
	}
	if fileCfg.MaxTokens != nil {
		c.MaxTokens = *fileCfg.MaxTokens
	}

	return nil
}
