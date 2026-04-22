package appconfig

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chenxiinn/go-llama-cpp-lab/pkg/chat"
)

const (
	LocalConfigPath        = "config/local.json"
	LocalExampleConfigPath = "config/local.example.json"
	UserConfigDirName      = ".go-llama-cpp-lab"
	UserConfigFileName     = "config.json"
)

type RuntimeBootstrap struct {
	ConfigPath string
	HomeDir    string
	Runtime    chat.RuntimeConfig
	WorkDir    string
}

func NewRuntimeBootstrap() RuntimeBootstrap {
	workDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	return RuntimeBootstrap{
		HomeDir: homeDir,
		Runtime: chat.DefaultRuntimeConfig(),
		WorkDir: workDir,
	}
}

func (b *RuntimeBootstrap) BindFlags(fs *flag.FlagSet) {
	fs.StringVar(&b.ConfigPath, "config", "", "Path to the config JSON file. If omitted, search ./config/local.json then ~/.go-llama-cpp-lab/config.json.")
	b.Runtime.BindFlags(fs)
}

func (b *RuntimeBootstrap) Load(args []string) (string, error) {
	explicitPath, explicit, err := ResolveExplicitConfigPath(args)
	if err != nil {
		return "", err
	}

	var path string
	if explicit {
		path = absolutePath(b.WorkDir, explicitPath)
		if err := ensureConfigExists(path); err != nil {
			return "", err
		}
	} else {
		path, err = FindDefaultConfigPath(b.WorkDir, b.HomeDir)
		if err != nil {
			return "", err
		}
	}

	if err := ApplyRuntimeConfigFile(path, &b.Runtime); err != nil {
		return "", err
	}

	b.ConfigPath = path
	return path, nil
}

func (b *RuntimeBootstrap) LoadAndParse(program string, args []string, bindExtra func(*flag.FlagSet)) (string, error) {
	loadedPath, err := b.Load(args)
	if err != nil {
		return "", err
	}

	fs := flag.NewFlagSet(program, flag.ExitOnError)
	b.BindFlags(fs)
	if bindExtra != nil {
		bindExtra(fs)
	}
	fs.Parse(args)

	if err := b.Runtime.Validate(); err != nil {
		return loadedPath, err
	}

	return loadedPath, nil
}

func ResolveExplicitConfigPath(args []string) (string, bool, error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-config" || arg == "--config":
			if i+1 >= len(args) {
				return "", false, errors.New("config flag requires a path value")
			}
			return args[i+1], true, nil
		case strings.HasPrefix(arg, "-config="):
			return strings.TrimPrefix(arg, "-config="), true, nil
		case strings.HasPrefix(arg, "--config="):
			return strings.TrimPrefix(arg, "--config="), true, nil
		}
	}

	return "", false, nil
}

func FindDefaultConfigPath(workDir, homeDir string) (string, error) {
	candidates := []string{
		absolutePath(workDir, LocalConfigPath),
		UserConfigPath(homeDir),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		exists, err := fileExists(candidate)
		if err != nil {
			return "", fmt.Errorf("stat config file %q: %w", candidate, err)
		}
		if exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no config file found; checked %q then %q", candidates[0], candidates[1])
}

func UserConfigPath(homeDir string) string {
	if homeDir == "" {
		return ""
	}
	return filepath.Join(homeDir, UserConfigDirName, UserConfigFileName)
}

func ApplyRuntimeConfigFile(path string, cfg *chat.RuntimeConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("decode config file %q: %w", path, err)
	}

	return nil
}

func absolutePath(baseDir, path string) string {
	if filepath.IsAbs(path) || baseDir == "" {
		return path
	}
	return filepath.Join(baseDir, path)
}

func ensureConfigExists(path string) error {
	exists, err := fileExists(path)
	if err != nil {
		return fmt.Errorf("stat config file %q: %w", path, err)
	}
	if !exists {
		return fmt.Errorf("config file %q does not exist", path)
	}
	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
