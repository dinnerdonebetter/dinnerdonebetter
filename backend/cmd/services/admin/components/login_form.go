package components

import (
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	"github.com/dinnerdonebetter/backend/internal/platform/webapp/loginform"

	g "maragu.dev/gomponents"
)

type ComponentRenderer struct {
	palette design.Palette
}

func NewComponentRenderer() *ComponentRenderer {
	return &ComponentRenderer{palette: design.StandardPalette}
}

// LoginFormProps is an alias for the shared login form props.
type LoginFormProps = loginform.Props

// LoginForm renders a login form using the shared component.
func (r *ComponentRenderer) LoginForm(props *LoginFormProps) g.Node {
	return loginform.Form(props, loginform.DefaultConfig(), &r.palette)
}
