package pages

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (b *Builder) AboutPage(ctx context.Context) gomponents.Node {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	now := time.Now()

	return components.PageShell("About",
		ghtml.H1(gomponents.Text("About")),

		ghtml.P(gomponents.Textf("Built with gomponents and rendered at %v.", now.Format(time.TimeOnly))),

		ghtml.P(
			gomponents.If(now.Second()%2 == 0, gomponents.Text("It's an even second!")),
			gomponents.If(now.Second()%2 != 0, gomponents.Text("It's an odd second!")),
		),

		ghtml.Img(
			ghtml.Class("max-w-sm"),
			ghtml.Src("https://www.gomponents.com/images/logo.png"),
			ghtml.Alt("gomponents logo"),
		),
	)
}
