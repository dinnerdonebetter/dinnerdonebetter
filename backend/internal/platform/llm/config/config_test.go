package llmcfg

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/llm/openai"
)

func TestConfig_ProvideLLMProvider_Empty(t *testing.T) {
	ctx := t.Context()
	cfg := &Config{Provider: ""}

	provider, err := cfg.ProvideLLMProvider(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider (noop)")
	}
}

func TestConfig_ProvideLLMProvider_OpenAI(t *testing.T) {
	ctx := t.Context()
	cfg := &Config{
		Provider: ProviderOpenAI,
		OpenAI: &openai.Config{
			APIKey: "test-key",
		},
	}

	provider, err := cfg.ProvideLLMProvider(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}
