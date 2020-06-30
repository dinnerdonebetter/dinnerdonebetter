package postgres

import "github.com/GuiaBolso/darwin"

func buildEquipmentMigration() darwin.Migration {
	return darwin.Migration{
		Version:     incrementMigrationVersion(),
		Description: "create valid instruments",
		Script: `
			INSERT INTO valid_instruments (
				"name",
				"variant",
				"description",
				"icon"
			) 
				VALUES
			(
				'knife', -- "name"
				'generic', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'knife', -- "name"
				'serrated', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'utility knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'steak knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'carving knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'carving fork', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cleaver', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'chinese cleaver', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'paring knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'boning knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'chef''s knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'santoku knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'bread knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cutting board', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pizza slicer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pizza stone', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'fat separator', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sifter', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sieve', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'can opener', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cups', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cup', -- "name"
				'1/8 cup', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cup', -- "name"
				'1/4 cup', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cup', -- "name"
				'1/3 cup', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cup', -- "name"
				'1/2 cup', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring cup', -- "name"
				'1 cup', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoons', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/8 teaspoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/4 teaspoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/2 teaspoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1 teaspoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/8 tablespoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/4 tablespoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1/2 tablespoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring spoon', -- "name"
				'1 tablespoon', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring bowl', -- "name"
				'large', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring bowl', -- "name"
				'medium', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'measuring bowl', -- "name"
				'small', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'colander', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'vegetable peeler', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'potato masher', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'potato ricer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'dough hook', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'meat tenderizer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'food mill', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'whisk', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'whisk', -- "name"
				'flat', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'salad spinner', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'grater', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'shears', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'citrus juicer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'juicer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'garlic press', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'14"', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'12"', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'10"', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'8"', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'skillet', -- "name"
				'cast iron', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'grill pan', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'grill pan', -- "name"
				'cast iron', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'saute pan', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'saute pan', -- "name"
				'6 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'saute pan', -- "name"
				'4 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'saute pan', -- "name"
				'3 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'saute pan', -- "name"
				'2 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'6 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'5 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'4 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'3 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'sauce pan', -- "name"
				'2 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'wok', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cassoulet', -- "name"
				'3 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'ramekin', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'steamer', -- "name"
				'8 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'steamer', -- "name"
				'5 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'steamer', -- "name"
				'3 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stockpot', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stockpot', -- "name"
				'8 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stockpot', -- "name"
				'12 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stockpot', -- "name"
				'16 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stockpot', -- "name"
				'32 quart', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'dutch oven', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'baking sheet', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'muffin pan', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'casserole dish', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'broiler pan', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'spatula', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'rubber/silicone spatula', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'scooper', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'scooper', -- "name"
				'large', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'ice scoop', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'spoon', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'fork', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'butter knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'slotted spoon', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'tongs', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'ladle', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'oven mitt', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'trivet', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'splatter guard', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'thermometer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'immersion blender', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'refrigerator', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stove', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'oven', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'scale', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'blender', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'towel', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'ice cube tray', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'mandoline', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'microplane', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'food processor', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'spice blender', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'mortar & pestle', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'molcajete', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'apple corer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'apple slicer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'baster', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'blow torch', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'bottle opener', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			), 
			(
				'mixing bowl', -- "name"
				'large', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'mixing bowl', -- "name"
				'medium', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'mixing bowl', -- "name"
				'small', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pie cutter', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cheese cutter', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cheese knife', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'conical sieve', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'cookie cutter', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'dough scraper', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'egg piercer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'egg separator', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'egg slicer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'egg timer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'fish scaler', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'fish spatula', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'funnel', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'honey dipper', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'lame', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'lobster pick', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'meat grinder', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pasta machine', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'melon baller', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'milk frother', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'nutcracker', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pastry bag', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pastry blender', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pastry brush', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pepper grinder', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'pot holder', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'rolling pin', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'roller docker', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'basket skimmer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'candy thermometer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'drum sieve', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'stand mixer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'hand mixer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'freezer', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'tea kettle', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			),
			(
				'zester', -- "name"
				'', -- "variant"
				'', -- "description"
				'' -- "icon"
			);`,
	}
}

func init() {
	migrations = append(
		migrations,
		buildEquipmentMigration(),
	)
}
