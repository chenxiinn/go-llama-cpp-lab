package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenxiinn/go-llama-cpp-lab/internal/appconfig"
)

type config struct {
	listenAddr string
	bootstrap  appconfig.RuntimeBootstrap
}

func newConfig() config {
	return config{
		listenAddr: "127.0.0.1:8080",
		bootstrap:  appconfig.NewRuntimeBootstrap(),
	}
}

func (c *config) bindExtraFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.listenAddr, "listen-addr", c.listenAddr, "Local HTTP listen address.")
}

func main() {
	cfg := newConfig()
	loadedPath, err := cfg.bootstrap.LoadAndParse(os.Args[0], os.Args[1:], cfg.bindExtraFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load server config: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "server placeholder: HTTP runtime lands in phase 6\n")
	fmt.Fprintf(os.Stdout, "config=%q listen_addr=%q model=%q n_ctx=%d gpu_layers=%d max_tokens=%d\n",
		loadedPath,
		cfg.listenAddr,
		cfg.bootstrap.Runtime.ModelPath,
		cfg.bootstrap.Runtime.ContextSize,
		cfg.bootstrap.Runtime.GPULayers,
		cfg.bootstrap.Runtime.MaxTokens,
	)
}
