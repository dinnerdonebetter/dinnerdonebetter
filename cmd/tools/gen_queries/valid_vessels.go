package main

const validVesselsTableName = "valid_vessels"

var validVesselsColumns = []string{
	idColumn,
	"name",
	"plural_name",
	"description",
	"icon_path",
	"usable_for_storage",
	"slug",
	"display_in_summary_lists",
	"include_in_generated_instructions",
	"capacity",
	"capacity_unit",
	"width_in_millimeters",
	"length_in_millimeters",
	"height_in_millimeters",
	"shape",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidVesselsQueries() []*Query {
	return []*Query{
		//
	}
}
