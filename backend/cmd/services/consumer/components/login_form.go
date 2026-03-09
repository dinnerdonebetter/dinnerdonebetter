package components

import (
	"github.com/dinnerdonebetter/backend/cmd/services/consumer/design"
	"github.com/dinnerdonebetter/backend/internal/webapp/loginform"

	g "maragu.dev/gomponents"
)

// LoginFormProps is an alias for the shared login form props.
type LoginFormProps = loginform.Props

// LoginForm renders a login form using the shared component.
func (r *ComponentRenderer) LoginForm(props *LoginFormProps) g.Node {
	palette := design.StandardPalette
	return loginform.Form(props, loginform.SignInConfig(), &palette)
}
