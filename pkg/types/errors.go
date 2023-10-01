package types

import (
	"errors"
)

var (
	errNewPasswordSameAsOld                 = errors.New("new password cannot be the same as the old password")
	errStartsAfterItEnds                    = errors.New("invalid start and end dates")
	errOneMainMinimumRequired               = errors.New("at least one main required for meal creation")
	errInvalidVotingDeadline                = errors.New("invalid voting deadline")
	errAtLeastOneRatingRequired             = errors.New("recipe rating must have at least one rating")
	errOneInstrumentOrVesselRequired        = errors.New("at least one instrument or vessel is required")
	errInstrumentIDOrProductIndicesRequired = errors.New("either instrumentID or productOfRecipeStepIndex and productOfRecipeStepProductIndex must be set")
	errDefaultValueMustBeEnumerationValue   = errors.New("default value must be in enumeration")
	errMustBeEitherMetricOrImperial         = errors.New("cannot be both metric and imperial")
	errInvalidType                          = errors.New("unexpected type received")
)
