package postgres

import "github.com/GuiaBolso/darwin"

func buildNonstandardProteinsMigration() darwin.Migration {
	return darwin.Migration{
		Version:     incrementMigrationVersion(),
		Description: "create valid instruments",
		Script: `
			INSERT INTO valid_ingredients (
				"name",
				"variant",
				"description",
				"warning",
				"contains_egg",
				"contains_dairy",
				"contains_peanut",
				"contains_tree_nut",
				"contains_soy",
				"contains_wheat",
				"contains_shellfish",
				"contains_sesame",
				"contains_fish",
				"contains_gluten",
				"animal_flesh",
				"animal_derived",
				"measurable_by_volume",
				"icon"
			)
				VALUES
			(
				'aluminum foil', -- "name"
				'', -- "variant"
				'', -- "description"
				'', -- "warning"
				'false', -- "contains_egg"
				'false', -- "contains_dairy"
				'false', -- "contains_peanut"
				'false', -- "contains_tree_nut"
				'false', -- "contains_soy"
				'false', -- "contains_wheat"
				'false', -- "contains_shellfish"
				'false', -- "contains_sesame"
				'false', -- "contains_fish"
				'false', -- "contains_gluten"
				'false', -- "animal_flesh"
				'false', -- "animal_derived"
				'true', -- "measurable_by_volume"
				'' -- "icon"
			),
			(
				'cheesecloth', -- "name"
				'', -- "variant"
				'', -- "description"
				'', -- "warning"
				'false', -- "contains_egg"
				'false', -- "contains_dairy"
				'false', -- "contains_peanut"
				'false', -- "contains_tree_nut"
				'false', -- "contains_soy"
				'false', -- "contains_wheat"
				'false', -- "contains_shellfish"
				'false', -- "contains_sesame"
				'false', -- "contains_fish"
				'false', -- "contains_gluten"
				'false', -- "animal_flesh"
				'false', -- "animal_derived"
				'false', -- "measurable_by_volume"
				'' -- "icon"
			),
			(
				'pastry bag', -- "name"
				'', -- "variant"
				'', -- "description"
				'', -- "warning"
				'false', -- "contains_egg"
				'false', -- "contains_dairy"
				'false', -- "contains_peanut"
				'false', -- "contains_tree_nut"
				'false', -- "contains_soy"
				'false', -- "contains_wheat"
				'false', -- "contains_shellfish"
				'false', -- "contains_sesame"
				'false', -- "contains_fish"
				'false', -- "contains_gluten"
				'false', -- "animal_flesh"
				'false', -- "animal_derived"
				'false', -- "measurable_by_volume"
				'' -- "icon"
			),
			(
				'parchment paper', -- "name"
				'', -- "variant"
				'', -- "description"
				'', -- "warning"
				'false', -- "contains_egg"
				'false', -- "contains_dairy"
				'false', -- "contains_peanut"
				'false', -- "contains_tree_nut"
				'false', -- "contains_soy"
				'false', -- "contains_wheat"
				'false', -- "contains_shellfish"
				'false', -- "contains_sesame"
				'false', -- "contains_fish"
				'false', -- "contains_gluten"
				'false', -- "animal_flesh"
				'false', -- "animal_derived"
				'true', -- "measurable_by_volume"
				'' -- "icon"
			);`,
	}
}

func init() {
	migrations = append(
		migrations,
		buildNonstandardProteinsMigration(),
	)
}
