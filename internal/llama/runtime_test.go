//go:build llama
// +build llama

package llama

import (
	"os"
	"testing"
)

const testModelEnv = "GO_LLAMA_CPP_LAB_TEST_MODEL"

func TestPhase2RuntimeSmoke(t *testing.T) {
	modelPath := os.Getenv(testModelEnv)
	if modelPath == "" {
		t.Skipf("set %s to run the Phase 2 model runtime smoke test", testModelEnv)
	}

	model, err := LoadModel(modelPath, DefaultModelConfig())
	if err != nil {
		t.Fatalf("load model: %v", err)
	}
	defer func() {
		if err := model.Close(); err != nil {
			t.Fatalf("close model: %v", err)
		}
	}()

	t.Run("ChatTemplate", func(t *testing.T) {
		tmpl, err := model.ChatTemplate()
		if err != nil {
			t.Fatalf("read chat template: %v", err)
		}
		if tmpl == "" {
			t.Fatal("chat template is empty")
		}
	})

	t.Run("TokenizeDecodeSample", func(t *testing.T) {
		prompt := "The capital of France is"

		tokens, err := model.Tokenize(prompt, true, true)
		if err != nil {
			t.Fatalf("tokenize prompt: %v", err)
		}
		if len(tokens) == 0 {
			t.Fatal("tokenize returned zero tokens")
		}

		ctx, err := model.NewContext(DefaultContextConfig())
		if err != nil {
			t.Fatalf("create context: %v", err)
		}
		defer func() {
			if err := ctx.Close(); err != nil {
				t.Fatalf("close context: %v", err)
			}
		}()

		samplerConfig := DefaultSamplerConfig()
		samplerConfig.UseGreedy = true
		sampler, err := NewSampler(samplerConfig)
		if err != nil {
			t.Fatalf("create sampler: %v", err)
		}
		defer func() {
			if err := sampler.Close(); err != nil {
				t.Fatalf("close sampler: %v", err)
			}
		}()

		if err := ctx.Decode(tokens); err != nil {
			t.Fatalf("decode prompt: %v", err)
		}
		if ctx.LogitsAt(-1) == nil {
			t.Fatal("last-token logits are nil after decode")
		}

		token, err := sampler.Sample(ctx)
		if err != nil {
			t.Fatalf("sample token: %v", err)
		}

		piece, err := model.TokenToPiece(token, 0, true)
		if err != nil {
			t.Fatalf("token to piece: %v", err)
		}
		if piece == "" {
			t.Fatalf("sampled token %d rendered to an empty piece", token)
		}

		if !model.IsEndOfGeneration(token) {
			if err := ctx.DecodeOne(token); err != nil {
				t.Fatalf("decode sampled token: %v", err)
			}
		}
	})
}
