package llama

import (
	"os"
	"path/filepath"
)

const (
	EnvLlamaCPPDir      = "LLAMA_CPP_DIR"
	EnvLlamaCPPBuildDir = "LLAMA_CPP_BUILD_DIR"
)

// LocalBuildPaths defines the local filesystem contract for the future cgo
// bridge. The actual bridge is added in phase 1; phase 0 only standardizes the
// directory inputs so build instructions and code use the same names.
type LocalBuildPaths struct {
	RootDir  string
	BuildDir string
}

// ResolveLocalBuildPaths reads the agreed environment variables for locating a
// local llama.cpp checkout and its build artifacts.
func ResolveLocalBuildPaths() LocalBuildPaths {
	rootDir := os.Getenv(EnvLlamaCPPDir)
	buildDir := os.Getenv(EnvLlamaCPPBuildDir)
	if buildDir == "" && rootDir != "" {
		buildDir = filepath.Join(rootDir, "build")
	}

	return LocalBuildPaths{
		RootDir:  rootDir,
		BuildDir: buildDir,
	}
}
