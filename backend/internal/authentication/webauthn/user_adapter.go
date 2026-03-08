package webauthn

import (
	"encoding/json"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnUser adapts identity.User and credentials to the webauthn.User interface.
type WebAuthnUser struct {
	User        *identity.User
	Credentials []*identity.WebAuthnCredential
	UserID      []byte // WebAuthn user handle - use user ID bytes
}

// Ensure WebAuthnUser implements webauthn.User.
var _ webauthn.User = (*WebAuthnUser)(nil)

// WebAuthnID returns the user handle (opaque byte sequence, max 64 bytes).
func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.UserID
}

// WebAuthnName returns the username for display during registration.
func (u *WebAuthnUser) WebAuthnName() string {
	return u.User.Username
}

// WebAuthnDisplayName returns the display name for registration.
func (u *WebAuthnUser) WebAuthnDisplayName() string {
	if u.User.FirstName != "" || u.User.LastName != "" {
		return strings.TrimSpace(u.User.FirstName + " " + u.User.LastName)
	}
	return u.User.Username
}

// WebAuthnCredentials returns credentials in webauthn.Credential format.
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	creds := make([]webauthn.Credential, 0, len(u.Credentials))
	for _, c := range u.Credentials {
		creds = append(creds, domainCredentialToWebAuthn(c))
	}
	return creds
}

func domainCredentialToWebAuthn(c *identity.WebAuthnCredential) webauthn.Credential {
	transports := parseTransports(c.Transports)
	return webauthn.Credential{
		ID:              c.CredentialID,
		PublicKey:       c.PublicKey,
		AttestationType: "none",
		Transport:       transports,
		Flags: webauthn.CredentialFlags{
			UserPresent:    true,
			UserVerified:   true,
			BackupEligible: false,
			BackupState:    false,
		},
		Authenticator: webauthn.Authenticator{
			SignCount:    c.SignCount,
			CloneWarning: false,
		},
	}
}

func parseTransports(s string) []protocol.AuthenticatorTransport {
	if s == "" {
		return nil
	}
	var transports []string
	if err := json.Unmarshal([]byte(s), &transports); err != nil {
		return nil
	}
	result := make([]protocol.AuthenticatorTransport, 0, len(transports))
	for _, t := range transports {
		result = append(result, protocol.AuthenticatorTransport(t))
	}
	return result
}
