package components

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"maragu.dev/gomponents"
	ghtmx "maragu.dev/gomponents-htmx"
	"maragu.dev/gomponents/html"
)

func columnHeader(name string) gomponents.Node {
	return html.Th(
		html.Class("border border-gray-300 px-4 py-2 text-left cursor-pointer"),
		ghtmx.Get("/sort?column="+name),
		gomponents.Text(name),
	)
}

func tableCell(content string) gomponents.Node {
	return html.Td(
		html.Class("border border-gray-300 px-4 py-2"),
		gomponents.Text(content),
	)
}

func buildColumnHeadersForType[T any]() gomponents.Node {
	var s T
	rowContents := []gomponents.Node{}
	for _, content := range GetFieldNames(s) {
		rowContents = append(rowContents, columnHeader(content))
	}

	return html.Tr(rowContents...)
}

func buildRowForType[T any](x T) gomponents.Node {
	rowContents := []gomponents.Node{}
	for _, content := range GetFieldValues(x) {
		rowContents = append(rowContents, tableCell(content))
	}

	return html.Tr(rowContents...)
}

func TableView[T any](newHref string, data *types.QueryFilteredResult[T]) gomponents.Node {
	rowComponents := []gomponents.Node{}
	for _, row := range data.Data {
		rowComponents = append(rowComponents, buildRowForType(*row))
	}

	descriptor := "items"
	if data.TotalCount == 1 {
		descriptor = "item"
	}

	return html.Div(
		html.Class("p-4 space-y-4"),
		html.Div(
			html.Class("flex justify-between items-center"),
			html.Div(
				html.Class("flex items-center space-x-2"),
				html.Input(
					html.Type("text"),
					html.Class("border border-gray-300 rounded-lg p-2 w-64"),
					html.Placeholder(fmt.Sprintf("Search from %d %s...", data.TotalCount, descriptor)),
				),
				html.Button(
					html.Class("bg-blue-500 text-white px-4 py-2 rounded-lg"),
					ghtmx.Boost("true"),
					html.Type("submit"),
					gomponents.Text("🔍"),
				),
			),
			html.A(
				html.Href(newHref),
				html.Class("bg-green-500 text-white px-4 py-2 rounded-lg"),
				gomponents.Text("New"),
			),
		),
		html.Div(
			html.Class("overflow-x-auto"),
			html.Table(
				html.Class("table-auto w-full border-collapse border border-gray-200"),
				html.THead(buildColumnHeadersForType[T]()),
				html.TBody(rowComponents...),
			),
		),
	)
}
