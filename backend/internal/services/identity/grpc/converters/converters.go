package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
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

// ConvertGRPCAccountInvitationUpdateRequestInputToAccountInvitationUpdateRequestInput creates a AccountInvitationDatabaseCreationInput from a AccountInvitationCreationRequestInput.
func ConvertGRPCAccountInvitationUpdateRequestInputToAccountInvitationUpdateRequestInput(input *identitysvc.AccountInvitationUpdateRequestInput) *identity.AccountInvitationUpdateRequestInput {
	x := &identity.AccountInvitationUpdateRequestInput{
		Token: input.Token,
		Note:  input.Note,
	}

	return x
}

func ConvertGRPCAccountCreationRequestInputToAccountCreationRequestInput(input *identitysvc.AccountCreationRequestInput) *identity.AccountCreationRequestInput {
	return &identity.AccountCreationRequestInput{
		Latitude:     pointer.To(float64(pointer.Dereference(input.Latitude))),
		Longitude:    pointer.To(float64(pointer.Dereference(input.Longitude))),
		Name:         input.Name,
		ContactPhone: input.ContactPhone,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		ZipCode:      input.ZipCode,
		Country:      input.Country,
	}
}

func ConvertAccountToGRPCAccount(input *identity.Account) *identitysvc.Account {
	var members []identitysvc.AccountUserMembershipWithUser
	for _, member := range input.Members {
		println(member) // TODO
	}

	return &identitysvc.Account{
		CreatedAt:                  grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		SubscriptionPlanID:         input.SubscriptionPlanID,
		LastUpdatedAt:              grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:                 grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Longitude:                  pointer.To(float32(pointer.Dereference(input.Longitude))),
		Latitude:                   pointer.To(float32(pointer.Dereference(input.Latitude))),
		State:                      input.State,
		ContactPhone:               input.ContactPhone,
		City:                       input.City,
		AddressLine1:               input.AddressLine1,
		ZipCode:                    input.ZipCode,
		Country:                    input.Country,
		BillingStatus:              input.BillingStatus,
		AddressLine2:               input.AddressLine2,
		PaymentProcessorCustomerID: input.PaymentProcessorCustomerID,
		BelongsToUser:              input.BelongsToUser,
		ID:                         input.ID,
		Name:                       input.Name,
		WebhookEncryptionKey:       input.WebhookEncryptionKey,
		Members:                    members,
	}
}
