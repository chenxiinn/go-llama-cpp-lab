//go:build llama
// +build llama

package llama

import "testing"

func TestProbeBackend(t *testing.T) {
	probe, err := ProbeBackend()
	if err != nil {
		t.Fatalf("probe backend: %v", err)
	}
	if probe.MaxDevices <= 0 {
		t.Fatalf("max devices = %d, want positive", probe.MaxDevices)
	}
	if probe.SystemInfo == "" {
		t.Fatal("system info is empty")
	}
}
