package postgres

import "github.com/GuiaBolso/darwin"

var (
	ingredientsMigration = darwin.Migration{
		Version:     1,
		Description: "create base ingredients",
		/* Sources:
		https://en.wikipedia.org/wiki/Edible_mushroom#Current_culinary_use
		https://en.wikipedia.org/wiki/List_of_edible_seeds
		https://en.wikipedia.org/wiki/List_of_culinary_nuts
		https://en.wikipedia.org/wiki/List_of_culinary_fruits
		https://en.wikipedia.org/wiki/List_of_root_vegetables
		https://en.wikipedia.org/wiki/List_of_vegetables
		https://en.wikipedia.org/wiki/List_of_leaf_vegetables
		https://en.wikipedia.org/wiki/List_of_dairy_products
		https://en.wikipedia.org/wiki/List_of_cheeses
		https://en.wikipedia.org/wiki/List_of_pasta
		https://en.wikipedia.org/wiki/Category:Staple_foods
		*/
		Script: `
			INSERT INTO TABLE ingredients
			(
				name,
				variant,
				description,
				warning,
				contains_egg,
				contains_dairy,
				contains_peanut,
				contains_tree_nut,
				contains_soy,
				contains_wheat,
				contains_shellfish,
				contains_sesame,
				contains_fish,
				contains_gluten,
				animal_flesh,
				animal_derived,
				considered_staple
			)
			VALUES
			(
				'sorghum', -- "name"
				'',        -- "variant",
				'',        -- "description",
				'',        -- "warning",
				'false',   -- "contains_egg",
				'false',   -- "contains_dairy",
				'false',   -- "contains_peanut",
				'false',   -- "contains_tree_nut",
				'false',   -- "contains_soy",
				'false',   -- "contains_wheat",
				'false',   -- "contains_shellfish",
				'false',   -- "contains_sesame",
				'false',   -- "contains_fish",
				'false',   -- "contains_gluten",
				'false',   -- "animal_flesh",
				'false',   -- "animal_derived",
				'false'    -- "considered_staple",
			),
			(
				'sorghum', -- "name"
				'',        -- "variant",
				'',        -- "description",
				'',        -- "warning",
				'false',   -- "contains_egg",
				'false',   -- "contains_dairy",
				'false',   -- "contains_peanut",
				'false',   -- "contains_tree_nut",
				'false',   -- "contains_soy",
				'false',   -- "contains_wheat",
				'false',   -- "contains_shellfish",
				'false',   -- "contains_sesame",
				'false',   -- "contains_fish",
				'false',   -- "contains_gluten",
				'false',   -- "animal_flesh",
				'false',   -- "animal_derived",
				'false'    -- "considered_staple",
			);
		`,
	}
)
