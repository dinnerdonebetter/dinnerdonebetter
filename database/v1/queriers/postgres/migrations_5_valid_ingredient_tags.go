package postgres

import "github.com/GuiaBolso/darwin"

func buildIngredientTagMigration() darwin.Migration {
	return darwin.Migration{
		Version:     incrementMigrationVersion(),
		Description: "create valid ingredient tags",
		Script: `
			INSERT INTO valid_ingredient_tags ( "name" ) VALUES
			( 'vegetarian' ),
			( 'vegan' ),
			( 'egg' ),
			( 'dairy' ),
			( 'peanut' ),
			( 'tree nut' ),
			( 'soy' ),
			( 'wheat' ),
			( 'shellfish' ),
			( 'sesame' ),
			( 'fish' ),
			( 'gluten' ),
			( 'animal flesh' ),
			( 'animal derived' ),
			( 'staple' );`,
	}
}

func init() {
	migrations = append(
		migrations,
		buildIngredientTagMigration(),
	)
}
