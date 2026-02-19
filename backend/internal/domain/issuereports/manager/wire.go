package manager

import "github.com/google/wire"

var (
	IssueReportsManagerProviders = wire.NewSet(
		NewIssueReportsDataManager,
	)
)
