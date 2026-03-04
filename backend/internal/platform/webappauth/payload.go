package webappauth

// AuthPayload is the struct stored in the secure cookie.
// Both admin and consumer apps use the same payload shape.
type AuthPayload struct {
	AccessToken string
}
