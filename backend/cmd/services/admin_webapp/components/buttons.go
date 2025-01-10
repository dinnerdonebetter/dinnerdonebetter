package components

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/components"
	ghtml "maragu.dev/gomponents/html"
)

func Button(text string) gomponents.Node {
	return ghtml.Button(
		combineComponentClasses(
			components.Classes{
				"w-full":              true,
				"px-2":                true,
				"py-1":                true,
				"font-semibold":       true,
				"transition-all":      true,
				"duration-300":        true,
				"rounded-lg":          true,
				"shadow-md":           true,
				"hover:shadow-lg":     true,
				"focus:outline-none":  true,
				"focus:ring-2":        true,
				"focus:ring-offset-2": true,
			},
			cssConstants.buttonClasses(),
		),
		gomponents.Text(text),
	)
}
