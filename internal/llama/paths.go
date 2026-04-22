package llama

import (
	"os"
	"path/filepath"
)

const (
	EnvLlamaCPPDir      = "LLAMA_CPP_DIR"
	EnvLlamaCPPBuildDir = "LLAMA_CPP_BUILD_DIR"
	EnvLlamaIncludeDir  = "LLAMA_INCLUDE_DIR"
	EnvLlamaLibDir      = "LLAMA_LIB_DIR"
)

// LocalBuildPaths defines the local filesystem contract for the future cgo
// bridge. Phase 1 uses these inputs to document and verify the native smoke
// path while still keeping the compile-time flags explicit in the build
// command.
type LocalBuildPaths struct {
	RootDir    string
	BuildDir   string
	IncludeDir string
	LibDir     string
}

// ResolveLocalBuildPaths reads the agreed environment variables for locating a
// local llama.cpp checkout and its build artifacts.
func ResolveLocalBuildPaths() LocalBuildPaths {
	rootDir := os.Getenv(EnvLlamaCPPDir)
	buildDir := os.Getenv(EnvLlamaCPPBuildDir)
	includeDir := os.Getenv(EnvLlamaIncludeDir)
	libDir := os.Getenv(EnvLlamaLibDir)
	if buildDir == "" && rootDir != "" {
		buildDir = filepath.Join(rootDir, "build")
	}
	if includeDir == "" && rootDir != "" {
		includeDir = filepath.Join(rootDir, "include")
	}
	if libDir == "" && buildDir != "" {
		libDir = filepath.Join(buildDir, "bin")
	}

	return LocalBuildPaths{
		RootDir:    rootDir,
		BuildDir:   buildDir,
		IncludeDir: includeDir,
		LibDir:     libDir,
	}
}
