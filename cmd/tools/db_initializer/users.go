package main

import (
	"context"
	"fmt"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// `PrixFixe1!` hashed (without backticks).
	defaultHashedPassword = "$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho"
	default2FASekret      = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
)

var (
	jonesHouseholdID = ""
)

var userCollection = struct {
	Admin,
	MomJones,
	DadJones,
	KidJones1,
	KidJones2 *types.UserDatabaseCreationInput
}{
	Admin: &types.UserDatabaseCreationInput{
		BirthMonth:           nil,
		BirthDay:             nil,
		ID:                   ksuid.New().String(),
		HashedPassword:       defaultHashedPassword,
		TwoFactorSecret:      default2FASekret,
		InvitationToken:      "",
		DestinationHousehold: "",
		Username:             "admin",
		EmailAddress:         "admin@prixfixe.email",
	},
	MomJones: &types.UserDatabaseCreationInput{
		BirthMonth:           nil,
		BirthDay:             nil,
		ID:                   ksuid.New().String(),
		HashedPassword:       defaultHashedPassword,
		TwoFactorSecret:      default2FASekret,
		InvitationToken:      "",
		DestinationHousehold: "",
		Username:             "momJones",
		EmailAddress:         "mom@jones.com",
	},
	DadJones: &types.UserDatabaseCreationInput{
		BirthMonth:           nil,
		BirthDay:             nil,
		ID:                   ksuid.New().String(),
		HashedPassword:       defaultHashedPassword,
		TwoFactorSecret:      default2FASekret,
		InvitationToken:      "",
		DestinationHousehold: "",
		Username:             "dadJones",
		EmailAddress:         "dad@jones.com",
	},
	KidJones1: &types.UserDatabaseCreationInput{
		BirthMonth:           nil,
		BirthDay:             nil,
		ID:                   ksuid.New().String(),
		HashedPassword:       defaultHashedPassword,
		TwoFactorSecret:      default2FASekret,
		InvitationToken:      "",
		DestinationHousehold: "",
		Username:             "sallyJones",
		EmailAddress:         "sally@jones.com",
	},
	KidJones2: &types.UserDatabaseCreationInput{
		BirthMonth:           nil,
		BirthDay:             nil,
		ID:                   ksuid.New().String(),
		HashedPassword:       defaultHashedPassword,
		TwoFactorSecret:      default2FASekret,
		InvitationToken:      "",
		DestinationHousehold: "",
		Username:             "billyJones",
		EmailAddress:         "billy@jones.com",
	},
}

func scaffoldUsers(ctx context.Context, db database.DataManager) error {
	if err := scaffoldJonesFamily(ctx, db); err != nil {
		return fmt.Errorf("creating jones family: %w", err)
	}

	if _, err := db.CreateUser(ctx, userCollection.Admin); err != nil {
		return fmt.Errorf("creating admin input: %w", err)
	}

	const adminUpdateQuery = `UPDATE users SET service_roles='service_admin', two_factor_secret_verified_on=extract(epoch FROM NOW()) WHERE username = 'admin';`
	if _, err := db.DB().ExecContext(ctx, adminUpdateQuery); err != nil {
		return fmt.Errorf("modifying admin input permissions: %w", err)
	}

	return nil
}

func scaffoldJonesFamily(ctx context.Context, db database.DataManager) error {
	momJones, err := db.CreateUser(ctx, userCollection.MomJones)
	if err != nil {
		return fmt.Errorf("creating mom input: %w", err)
	}

	jonesHouseholdID, err = db.GetDefaultHouseholdIDForUser(ctx, userCollection.MomJones.ID)
	if err != nil {
		return fmt.Errorf("fetching default household for mom input: %w", err)
	}

	type userContainer struct {
		input  *types.UserDatabaseCreationInput
		invite *types.HouseholdInvitation
	}

	dadsInviteInput := &types.HouseholdInvitationDatabaseCreationInput{
		ID:                   ksuid.New().String(),
		FromUser:             momJones.ID,
		Note:                 "",
		ToEmail:              userCollection.DadJones.EmailAddress,
		Token:                "example_invite_token_dad",
		DestinationHousehold: jonesHouseholdID,
	}

	dadsInvite, err := db.CreateHouseholdInvitation(ctx, dadsInviteInput)
	if err != nil {
		return fmt.Errorf("creating dad's invite: %w", err)
	}

	sallysInviteInput := &types.HouseholdInvitationDatabaseCreationInput{
		ID:                   ksuid.New().String(),
		FromUser:             momJones.ID,
		Note:                 "",
		ToEmail:              userCollection.KidJones1.EmailAddress,
		Token:                "example_invite_token_sally",
		DestinationHousehold: jonesHouseholdID,
	}

	sallysInvite, err := db.CreateHouseholdInvitation(ctx, sallysInviteInput)
	if err != nil {
		return fmt.Errorf("creating sally's invite: %w", err)
	}

	billysInviteInput := &types.HouseholdInvitationDatabaseCreationInput{
		ID:                   ksuid.New().String(),
		FromUser:             momJones.ID,
		Note:                 "",
		ToEmail:              userCollection.KidJones2.EmailAddress,
		Token:                "example_invite_token_billy",
		DestinationHousehold: jonesHouseholdID,
	}

	billysInvite, err := db.CreateHouseholdInvitation(ctx, billysInviteInput)
	if err != nil {
		return fmt.Errorf("creating billy's invite: %w", err)
	}

	jonesUsers := map[string]userContainer{
		userCollection.DadJones.Username: {
			input:  userCollection.DadJones,
			invite: dadsInvite,
		},
		userCollection.KidJones1.Username: {
			input:  userCollection.KidJones1,
			invite: sallysInvite,
		},
		userCollection.KidJones2.Username: {
			input:  userCollection.KidJones2,
			invite: billysInvite,
		},
	}

	for username, cfg := range jonesUsers {
		cfg.input.InvitationToken = cfg.invite.Token
		cfg.input.DestinationHousehold = cfg.invite.DestinationHousehold
		if _, err = db.CreateUser(ctx, cfg.input); err != nil {
			return fmt.Errorf("creating Jones family input %q: %w", username, err)
		}
	}

	return nil
}
