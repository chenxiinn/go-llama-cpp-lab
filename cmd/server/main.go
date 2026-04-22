package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenxiinn/go-llama-cpp-lab/pkg/chat"
)

type config struct {
	listenAddr string
	runtime    chat.RuntimeConfig
}

func newConfig() config {
	return config{
		listenAddr: "127.0.0.1:8080",
		runtime:    chat.DefaultRuntimeConfig(),
	}
}

func (c *config) bindFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.listenAddr, "listen-addr", c.listenAddr, "Local HTTP listen address.")
	c.runtime.BindFlags(fs)
}

func main() {
	cfg := newConfig()

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg.bindFlags(fs)
	fs.Parse(os.Args[1:])

	if err := cfg.runtime.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid server config: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "server placeholder: HTTP runtime lands in phase 6\n")
	fmt.Fprintf(os.Stdout, "listen_addr=%q model=%q n_ctx=%d gpu_layers=%d max_tokens=%d\n",
		cfg.listenAddr,
		cfg.runtime.ModelPath,
		cfg.runtime.ContextSize,
		cfg.runtime.GPULayers,
		cfg.runtime.MaxTokens,
	)
}
