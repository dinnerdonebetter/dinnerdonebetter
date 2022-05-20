package main

import (
	"context"
	"fmt"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"

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
	var err error
	if _, err = db.CreateUser(ctx, userCollection.MomJones); err != nil {
		return fmt.Errorf("creating mom user: %w", err)
	}

	jonesHouseholdID, err = db.GetDefaultHouseholdIDForUser(ctx, userCollection.MomJones.ID)
	if err != nil {
		return fmt.Errorf("fetching default household for mom user: %w", err)
	}

	jonesUsers := []*types.UserDatabaseCreationInput{
		userCollection.DadJones,
		userCollection.KidJones1,
		userCollection.KidJones2,
	}

	for _, input := range jonesUsers {
		input.DestinationHousehold = jonesHouseholdID
		if _, err = db.CreateUser(ctx, input); err != nil {
			return fmt.Errorf("creating Jones family user %q: %w", input.Username, err)
		}
	}

	if _, err = db.CreateUser(ctx, userCollection.Admin); err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}

	adminHouseholdID, err := db.GetDefaultHouseholdIDForUser(ctx, userCollection.Admin.ID)
	if err != nil {
		return fmt.Errorf("fetching default household for admin user: %w", err)
	}

	err = db.ModifyUserPermissions(ctx, adminHouseholdID, userCollection.Admin.ID, &types.ModifyUserPermissionsInput{
		Reason: "scaffold",
		NewRoles: []string{
			authorization.ServiceAdminRole.String(),
			authorization.ServiceUserRole.String(),
		},
	})
	if err != nil {
		return fmt.Errorf("modifying admin user permissions: %w", err)
	}

	return nil
}
