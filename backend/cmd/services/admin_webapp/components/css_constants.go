package components

import (
	"fmt"

	"maragu.dev/gomponents/components"
)

type cssConstantManger struct {
	PrimaryColor    string
	ButtonTextColor string
}

var (
	cssConstants = cssConstantManger{
		PrimaryColor:    "blue",
		ButtonTextColor: "white",
	}
)

func (m *cssConstantManger) buttonClasses() components.Classes {
	return components.Classes{
		fmt.Sprintf("bg-%s-600", m.PrimaryColor):         true,
		fmt.Sprintf("hover:bg-%s-700", m.PrimaryColor):   true,
		fmt.Sprintf("focus:ring-%s-400", m.PrimaryColor): true,
		fmt.Sprintf("text-%s", m.ButtonTextColor):        true,
	}
}

func combineComponentClasses(classes ...components.Classes) components.Classes {
	combined := components.Classes{}
	for _, c := range classes {
		for k, v := range c {
			combined[k] = v
		}
	}
	return combined
}
