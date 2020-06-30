package postgres

import "github.com/GuiaBolso/darwin"

func buildPreparationsMigration() darwin.Migration {
	return darwin.Migration{
		Version:     incrementMigrationVersion(),
		Description: "create valid preparations",
		Script: `
			INSERT INTO valid_preparations (
				"name",
				"description",
				"icon",
				"applicable_to_all_ingredients"
			) 
				VALUES
			(
				'bake', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'boil', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'blanch', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'braise', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'coddle', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'infuse', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'pressure cook', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'simmer', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'poach', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'steam', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'double steam', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'steep', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'stew', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'grill', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'barbecue', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'fry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'deep fry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'pan fry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'sauté', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'stir fry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'shallow fry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'microwave', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'roast', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'smoke', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'sear', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'brine', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'dry', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'ferment', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'marinate', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'pickle', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'season', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'sour', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'sprout', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'cut', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'slice', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'dice', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'grate', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'julienne', -- "name"
				'measures approximately 1⁄8 by 1⁄8 by 1–2 inches (0.3 cm × 0.3 cm × 3 cm–5 cm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'fine julienne', -- "name"
				'measures approximately 1⁄16 by 1⁄16 by 1–2 inches', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'pont-neuf', -- "name"
				'measures from 1⁄3 by 1⁄3 by 2 1⁄2–3 inches (1 cm × 1 cm × 6 cm–8 cm) to 3⁄4 by 3⁄4 by 3 inches (2 cm × 2 cm × 8 cm).', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'batonnet', -- "name"
				'measures approximately 1⁄4 by 1⁄4 by 2–2 1⁄2 inches (0.6 cm × 0.6 cm × 5 cm–6 cm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'chiffonade', -- "name"
				'To roll and cut in sections from 4-10mm in width', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'large dice', -- "name"
				'sides measuring approximately 3⁄4 inch (20 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'medium dice', -- "name"
				'sides measuring approximately 1⁄2 inch', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'small dice', -- "name"
				'sides measuring approximately 1⁄4 inch (5 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'brunoise', -- "name"
				'sides measuring approximately 1⁄8 inch (3 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'fine brunoise', -- "name"
				'sides measuring approximately 1⁄16 inch (2 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'paysanne', -- "name"
				'1⁄2 by 1⁄2 by 1⁄8 inch (10 mm × 10 mm × 3 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'lozenge', -- "name"
				'diamond shape, 1⁄2 by 1⁄2 by 1⁄8 inch (10 mm × 10 mm × 3 mm)', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'tourné', -- "name"
				'2 inches (50 mm) long with seven faces usually with a bulge in the center portion', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'vacuum seal', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'mince', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'peel', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'shave', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'knead', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'mill', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'grind', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'mix', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'blend', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			),
			(
				'', -- "name"
				'', -- "description"
				'', -- "icon"
				'false' -- "applicable_to_all_ingredients"
			);`,
	}
}

func init() {
	migrations = append(
		migrations,
		buildPreparationsMigration(),
	)
}
