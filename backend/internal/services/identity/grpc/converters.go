package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
)

func ConvertGRPCUserDetailsUpdateRequestInputToUserDetailsDatabaseUpdateInput(input *identitysvc.UserDetailsUpdateRequestInput) *identity.UserDetailsDatabaseUpdateInput {
	return &identity.UserDetailsDatabaseUpdateInput{
		Birthday:  grpcconverters.ConvertPBTimestampToTime(input.Birthday),
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
}

func ConvertUserToGRPCUser(input *identity.User) *identitysvc.User {
	return &identitysvc.User{
		CreatedAt:                  grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		PasswordLastChangedAt:      grpcconverters.ConvertTimePointerToPBTimestamp(input.PasswordLastChangedAt),
		LastUpdatedAt:              grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		LastAcceptedTermsOfService: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  grpcconverters.ConvertTimePointerToPBTimestamp(input.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  grpcconverters.ConvertTimePointerToPBTimestamp(input.TwoFactorSecretVerifiedAt),
		EmailAddressVerifiedAt:     grpcconverters.ConvertTimePointerToPBTimestamp(input.EmailAddressVerifiedAt),
		Birthday:                   grpcconverters.ConvertTimePointerToPBTimestamp(input.Birthday),
		ArchivedAt:                 grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		HashedPassword:             input.HashedPassword,
		LastName:                   input.LastName,
		AccountStatusExplanation:   input.AccountStatusExplanation,
		ID:                         input.ID,
		AccountStatus:              input.AccountStatus,
		Username:                   input.Username,
		FirstName:                  input.FirstName,
		TwoFactorSecret:            input.TwoFactorSecret,
		EmailAddress:               input.EmailAddress,
		AvatarSrc:                  input.AvatarSrc,
		ServiceRole:                input.ServiceRole,
		RequiresPasswordChange:     input.RequiresPasswordChange,
	}
}
