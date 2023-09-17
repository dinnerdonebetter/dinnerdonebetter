package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

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
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validVesselsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validVesselsTableName,
				strings.Join(validVesselsColumns, ",\n\t"),
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidVesselExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validVesselsTableName,
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				validVesselsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVesselIDsNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRandomValidVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVesselsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidVesselLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validVesselsTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
