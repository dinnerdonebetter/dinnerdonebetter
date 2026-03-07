package llm

import (
	"context"
	"testing"
)

func TestNoopProvider_Completion(t *testing.T) {
	ctx := context.Background()
	provider := NewNoopProvider()

	result, err := provider.Completion(ctx, CompletionParams{
		Model: "test",
		Messages: []Message{
			{Role: "user", Content: "hello"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Content != "" {
		t.Errorf("expected empty content, got %q", result.Content)
	}
}
