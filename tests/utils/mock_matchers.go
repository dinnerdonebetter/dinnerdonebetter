package testutils

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

// ContextMatcher is a matcher for use with testify/mock's MatchBy function. It provides some level of type
// safety reassurance over mock.Anything, in that the resulting function will panic if anything other than
// a context.Context.
var ContextMatcher interface{} = mock.MatchedBy(func(context.Context) bool {
	return true
})

// HTTPRequestMatcher is a matcher. It provides some level of type assurance over mock.Anything,
// in that the resulting function will panic if anything other than a *http.Request.
var HTTPRequestMatcher interface{} = mock.MatchedBy(func(*http.Request) bool {
	return true
})

// HTTPResponseWriterMatcher is a matcher for the http.ResponseWriter interface. It provides some level of type
// safety reassurance over mock.Anything, in that the resulting function will panic if anything other than
// a http.ResponseWriter.
var HTTPResponseWriterMatcher interface{} = mock.MatchedBy(func(http.ResponseWriter) bool {
	return true
})

// BuildAuditLogEntryCreationInputMatcher is a matcher for use with testify/mock's MatchBy function.
func BuildAuditLogEntryCreationInputMatcher(event *types.AuditLogEntryCreationInput) func(*types.AuditLogEntryCreationInput) bool {
	return func(input *types.AuditLogEntryCreationInput) bool {
		for k, v := range input.Context {
			if v2, present := event.Context[k]; !present || v2 != v {
				return false
			}
		}

		return input.EventType == event.EventType
	}
}

// BuildAuditLogEntryCreationInputEventTypeMatcher is a matcher for use with testify/mock's MatchBy function.
func BuildAuditLogEntryCreationInputEventTypeMatcher(eventType string) func(*types.AuditLogEntryCreationInput) bool {
	return func(input *types.AuditLogEntryCreationInput) bool {
		return input.EventType == eventType
	}
}
