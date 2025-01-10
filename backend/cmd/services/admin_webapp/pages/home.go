package pages

import (
	"context"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (b *PageBuilder) HomePage(ctx context.Context) gomponents.Node {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return components.PageShell("Home",
		ghtml.H1(gomponents.Text("Home")),
		ghtml.P(gomponents.Text("welcome.")),
	)
}
