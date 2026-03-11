package queuetest

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	queuetest "github.com/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"
)

// BuildResult holds the queue test job and a cleanup that flushes metrics so
// short-lived CronJob pods export queue_test_round_trip_ms before exit.
type BuildResult struct {
	Job             *queuetest.Job
	ShutdownMetrics func(context.Context) error
}

// NewBuildResult constructs a BuildResult so wire can inject both job and metrics shutdown.
func NewBuildResult(job *queuetest.Job, provider metrics.Provider) *BuildResult {
	return &BuildResult{Job: job, ShutdownMetrics: provider.Shutdown}
}
