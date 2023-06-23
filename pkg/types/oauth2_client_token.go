package types

import (
	"context"
	"time"
)

type (
	// OAuth2ClientToken represents a user-authorized OAuth2 client's token.
	OAuth2ClientToken struct {
		_                   struct{}
		RefreshCreateAt     time.Time     `json:"refreshCreateAt"`
		AccessCreateAt      time.Time     `json:"accessCreateAt"`
		CodeCreateAt        time.Time     `json:"codeCreateAt"`
		RedirectURI         string        `json:"redirectURI"`
		Scope               string        `json:"scope"`
		Code                string        `json:"code"`
		CodeChallenge       string        `json:"codeChallenge"`
		CodeChallengeMethod string        `json:"codeChallengeMethod"`
		BelongsToUser       string        `json:"belongsToUser"`
		Access              string        `json:"access"`
		ClientID            string        `json:"clientID"`
		Refresh             string        `json:"refresh"`
		ID                  string        `json:"id"`
		CodeExpiresIn       time.Duration `json:"codeExpiresIn"`
		AccessExpiresIn     time.Duration `json:"accessExpiresIn"`
		RefreshExpiresIn    time.Duration `json:"refreshExpiresIn"`
	}

	// OAuth2ClientTokenDatabaseCreationInput represents a user-authorized OAuth2 client's token's database creation input.
	OAuth2ClientTokenDatabaseCreationInput struct {
		_                   struct{}
		RefreshCreateAt     time.Time
		AccessCreateAt      time.Time
		CodeCreateAt        time.Time
		RedirectURI         string
		Scope               string
		Code                string
		CodeChallenge       string
		CodeChallengeMethod string
		BelongsToUser       string
		Access              string
		ClientID            string
		Refresh             string
		ID                  string
		CodeExpiresIn       time.Duration
		AccessExpiresIn     time.Duration
		RefreshExpiresIn    time.Duration
	}

	OAuth2ClientTokenDataManager interface {
		GetOAuth2ClientTokenByCode(ctx context.Context, code string) (*OAuth2ClientToken, error)
		GetOAuth2ClientTokenByAccess(ctx context.Context, access string) (*OAuth2ClientToken, error)
		GetOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) (*OAuth2ClientToken, error)
		CreateOAuth2ClientToken(ctx context.Context, input *OAuth2ClientTokenDatabaseCreationInput) (*OAuth2ClientToken, error)
		ArchiveOAuth2ClientTokenByAccess(ctx context.Context, access string) error
		ArchiveOAuth2ClientTokenByCode(ctx context.Context, code string) error
		ArchiveOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) error
	}
)
