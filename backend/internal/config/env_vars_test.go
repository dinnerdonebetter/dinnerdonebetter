package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCeaseOperation(T *testing.T) {
	T.Run("unset", func(t *testing.T) {
		assert.False(t, ShouldCeaseOperation())
	})

	T.Run("set", func(t *testing.T) {
		t.Setenv(CeaseOperationEnvVarKey, "true")

		assert.True(t, ShouldCeaseOperation())
	})
}

func TestRunningInTheCloud(T *testing.T) {
	T.Run("unset", func(t *testing.T) {
		assert.False(t, RunningInTheCloud())
	})

	T.Run("set", func(t *testing.T) {
		t.Setenv(RunningInGCPEnvVarKey, "true")

		assert.True(t, RunningInTheCloud())
	})
}

func TestRunningInKubernetes(T *testing.T) {
	T.Run("unset", func(t *testing.T) {
		assert.False(t, RunningInKubernetes())
	})

	T.Run("set", func(t *testing.T) {
		t.Setenv(RunningInKubernetesEnvVarKey, "true")

		assert.True(t, RunningInKubernetes())
	})
}
