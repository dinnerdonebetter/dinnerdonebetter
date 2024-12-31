package config

import (
	"os"
	"strings"
)

const (
	// CeaseOperationEnvVarKey is the env var key used to indicate a function or job should just quit early.
	CeaseOperationEnvVarKey = "CEASE_OPERATION"
	// RunningInGCPEnvVarKey is the env var key we use to indicate we're running in GCP.
	RunningInGCPEnvVarKey = "RUNNING_IN_GCP"
)

func ShouldCeaseOperation() bool {
	return strings.TrimSpace(strings.ToLower(os.Getenv(CeaseOperationEnvVarKey))) == "true"
}

func RunningInTheCloud() bool {
	return os.Getenv(RunningInGCPEnvVarKey) != ""
}
