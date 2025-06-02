package components

import (
	"fmt"

	"maragu.dev/gomponents"
	ghtmx "maragu.dev/gomponents-htmx"
	ghtml "maragu.dev/gomponents/html"
)

type SubmissionFormProps struct {
	PostAddress string
	TargetID    string
}

func BuildHTMXPoweredSubmissionForm(props SubmissionFormProps, formContents ...gomponents.Node) gomponents.Node {
	return ghtml.Form(
		ghtml.Class("space-y-4"),
		ghtmx.Post(props.PostAddress),
		ghtmx.Ext("json-enc"),
		ghtmx.Target(fmt.Sprintf("#%s", props.TargetID)),
		ghtmx.Swap("innerHTML"),
		ghtml.Div(formContents...),
		Button("Submit"),
	)
}
