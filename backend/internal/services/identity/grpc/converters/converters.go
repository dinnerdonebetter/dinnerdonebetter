package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func ConvertGRPCUserDetailsUpdateRequestInputToUserDetailsUpdateRequestInput(input *identitysvc.UserDetailsUpdateRequestInput) *identity.UserDetailsUpdateRequestInput {
	return &identity.UserDetailsUpdateRequestInput{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Birthday:        grpcconverters.ConvertPBTimestampToTime(input.Birthday),
		CurrentPassword: input.CurrentPassword,
		TOTPToken:       input.TOTPToken,
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

func ConvertAccountCreationRequestInputToGRPCAccountCreationRequestInput(input *identity.AccountCreationRequestInput) *identitysvc.AccountCreationRequestInput {
	return &identitysvc.AccountCreationRequestInput{
		Latitude:     pointer.To(float32(pointer.Dereference(input.Latitude))),
		Longitude:    pointer.To(float32(pointer.Dereference(input.Longitude))),
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
	var members []*identitysvc.AccountUserMembershipWithUser
	for _, member := range input.Members {
		members = append(members, ConvertAccountUserMembershipWithUserToGRPCAccountUserMembershipWithUser(member))
	}

	return &identitysvc.Account{
		CreatedAt:                  grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:              grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:                 grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Longitude:                  pointer.To(float32(pointer.Dereference(input.Longitude))),
		Latitude:                   pointer.To(float32(pointer.Dereference(input.Latitude))),
		SubscriptionPlanID:         input.SubscriptionPlanID,
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

func ConvertAccountUserMembershipWithUserToGRPCAccountUserMembershipWithUser(input *identity.AccountUserMembershipWithUser) *identitysvc.AccountUserMembershipWithUser {
	return &identitysvc.AccountUserMembershipWithUser{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		BelongsToUser:    ConvertUserToGRPCUser(input.BelongsToUser),
		ID:               input.ID,
		BelongsToAccount: input.BelongsToAccount,
		AccountRole:      input.AccountRole,
		DefaultAccount:   input.DefaultAccount,
	}
}

func ConvertGRPCAccountInvitationCreationRequestInputToAccountInvitationCreationRequestInput(input *identitysvc.AccountInvitationCreationRequestInput) *identity.AccountInvitationCreationRequestInput {
	return &identity.AccountInvitationCreationRequestInput{
		ExpiresAt: grpcconverters.ConvertPBTimestampToTimePointer(input.ExpiresAt),
		Note:      input.Note,
		ToEmail:   input.ToEmail,
		ToName:    input.ToName,
	}
}

func ConvertAccountInvitationToGRPCAccountInvitation(input *identity.AccountInvitation) *identitysvc.AccountInvitation {
	return &identitysvc.AccountInvitation{
		CreatedAt:          grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:      grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:         grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		ToUser:             input.ToUser,
		Status:             input.Status,
		ToEmail:            input.ToEmail,
		StatusNote:         input.StatusNote,
		Token:              input.Token,
		ID:                 input.ID,
		Note:               input.Note,
		ToName:             input.ToName,
		ExpiresAt:          grpcconverters.ConvertTimeToPBTimestamp(input.ExpiresAt),
		DestinationAccount: ConvertAccountToGRPCAccount(&input.DestinationAccount),
		FromUser:           ConvertUserToGRPCUser(&input.FromUser),
	}
}

func ConvertGRPCUserRegistrationInputToUserRegistrationInput(input *identitysvc.UserRegistrationInput) *identity.UserRegistrationInput {
	return &identity.UserRegistrationInput{
		Birthday:              grpcconverters.ConvertPBTimestampToTimePointer(input.Birthday),
		Password:              input.Password,
		EmailAddress:          input.EmailAddress,
		InvitationToken:       input.InvitationToken,
		InvitationID:          input.InvitationID,
		Username:              input.Username,
		FirstName:             input.FirstName,
		LastName:              input.LastName,
		AccountName:           input.AccountName,
		AcceptedTOS:           input.AcceptedTOS,
		AcceptedPrivacyPolicy: input.AcceptedPrivacyPolicy,
	}
}

func ConvertUserRegistrationInputToGRPCUserRegistrationInput(input *identity.UserRegistrationInput) *identitysvc.UserRegistrationInput {
	return &identitysvc.UserRegistrationInput{
		Birthday:              grpcconverters.ConvertTimePointerToPBTimestamp(input.Birthday),
		Password:              input.Password,
		EmailAddress:          input.EmailAddress,
		InvitationToken:       input.InvitationToken,
		InvitationID:          input.InvitationID,
		Username:              input.Username,
		FirstName:             input.FirstName,
		LastName:              input.LastName,
		AccountName:           input.AccountName,
		AcceptedTOS:           input.AcceptedTOS,
		AcceptedPrivacyPolicy: input.AcceptedPrivacyPolicy,
	}
}

func ConvertUserCreationResponseToGRPCUserCreationResponse(input *identity.UserCreationResponse) *identitysvc.UserCreationResponse {
	return &identitysvc.UserCreationResponse{
		CreatedAt:       grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		Birthday:        grpcconverters.ConvertTimePointerToPBTimestamp(input.Birthday),
		Username:        input.Username,
		EmailAddress:    input.EmailAddress,
		TwoFactorQRCode: input.TwoFactorQRCode,
		CreatedUserID:   input.CreatedUserID,
		AccountStatus:   input.AccountStatus,
		TwoFactorSecret: input.TwoFactorSecret,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
	}
}

func ConvertGRPCAccountOwnershipTransferInputToAccountOwnershipTransferInput(input *identitysvc.AccountOwnershipTransferInput) *identity.AccountOwnershipTransferInput {
	return &identity.AccountOwnershipTransferInput{
		Reason:       input.Reason,
		CurrentOwner: input.CurrentOwner,
		NewOwner:     input.NewOwner,
	}
}

func convertFloat32PointerToFloat64Pointer(input *float32) *float64 {
	if input == nil {
		return nil
	}
	return pointer.To(float64(*input))
}

func ConvertGRPCAccountUpdateRequestInputToAccountUpdateRequestInput(input *identitysvc.AccountUpdateRequestInput) *identity.AccountUpdateRequestInput {
	return &identity.AccountUpdateRequestInput{
		Name:          input.Name,
		ContactPhone:  input.ContactPhone,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		City:          input.City,
		State:         input.State,
		ZipCode:       input.ZipCode,
		Country:       input.Country,
		Latitude:      convertFloat32PointerToFloat64Pointer(input.Latitude),
		Longitude:     convertFloat32PointerToFloat64Pointer(input.Longitude),
		BelongsToUser: input.BelongsToUser,
	}
}

func ConvertGRPCModifyUserPermissionsInputToModifyUserPermissionsInput(input *identitysvc.ModifyUserPermissionsInput) *identity.ModifyUserPermissionsInput {
	return &identity.ModifyUserPermissionsInput{
		Reason:  input.Reason,
		NewRole: input.NewRole,
	}
}
