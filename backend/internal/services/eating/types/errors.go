package types

import (
	"errors"
)

var (
	ErrNewPasswordSameAsOld                 = errors.New("new password cannot be the same as the old password")
	ErrStartsAfterItEnds                    = errors.New("invalid start and end dates")
	ErrOneMainMinimumRequired               = errors.New("at least one main required for meal creation")
	ErrInvalidVotingDeadline                = errors.New("invalid voting deadline")
	ErrAtLeastOneRatingRequired             = errors.New("recipe rating must have at least one rating")
	ErrOneInstrumentOrVesselRequired        = errors.New("at least one instrument or vessel is required")
	ErrInstrumentIDOrProductIndicesRequired = errors.New("either instrumentID or productOfRecipeStepIndex and productOfRecipeStepProductIndex must be set")
	ErrDefaultValueMustBeEnumerationValue   = errors.New("default value must be in enumeration")
	ErrMustBeEitherMetricOrImperial         = errors.New("cannot be both metric and imperial")
	ErrInvalidType                          = errors.New("unexpected type received")

	// this just ensures that we don't have any duplicated codes.
	_ = map[string]ErrorCode{
		string(ErrFetchingSessionContextData): ErrFetchingSessionContextData,
		string(ErrDecodingRequestInput):       ErrDecodingRequestInput,
		string(ErrValidatingRequestInput):     ErrValidatingRequestInput,
		string(ErrDataNotFound):               ErrDataNotFound,
		string(ErrTalkingToDatabase):          ErrTalkingToDatabase,
		string(ErrTalkingToSearchProvider):    ErrTalkingToSearchProvider,
	}
)

type (
	errorCode string
	ErrorCode errorCode
)

const (
	// ErrNothingSpecific is a catch-all error code for when we just need one.
	ErrNothingSpecific ErrorCode = "E100"
	// ErrFetchingSessionContextData is returned when we fail to fetch session context data.
	ErrFetchingSessionContextData ErrorCode = "E101"
	// ErrDecodingRequestInput is returned when we fail to decode request input.
	ErrDecodingRequestInput ErrorCode = "E102"
	// ErrValidatingRequestInput is returned when the user provides invalid input.
	ErrValidatingRequestInput ErrorCode = "E103"
	// ErrDataNotFound is returned when we fail to find data in the database.
	ErrDataNotFound ErrorCode = "E104"
	// ErrTalkingToDatabase is returned when we fail to interact with a database.
	ErrTalkingToDatabase ErrorCode = "E105"
	// ErrMisbehavingDependency is returned when we fail to interact with a third party.
	ErrMisbehavingDependency ErrorCode = "E106"
	// ErrTalkingToSearchProvider is returned when we fail to interact with the search provider.
	ErrTalkingToSearchProvider ErrorCode = "E107"
	// ErrSecretGeneration is returned when a user is not authorized.
	ErrSecretGeneration ErrorCode = "E108"
	// ErrUserIsBanned is returned when a user is banned.
	ErrUserIsBanned ErrorCode = "E109"
	// ErrUserIsNotAuthorized is returned when a user is not authorized.
	ErrUserIsNotAuthorized ErrorCode = "E110"
	// ErrEncryptionIssue is returned when encryption fails in the service.
	ErrEncryptionIssue ErrorCode = "E111"
)
