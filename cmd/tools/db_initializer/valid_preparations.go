package main

import "gitlab.com/prixfixe/prixfixe/pkg/types"

var (
	validPreparations = []*types.ValidPreparationCreationInput{
		{
			Name:        "bake",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "boil",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "blanch",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "braise",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "coddle",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "infuse",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "pressure cook",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "simmer",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "poach",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "steam",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "double steam",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "steep",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "stew",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "grill",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "barbecue",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "fry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "deep fry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "pan fry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "sauté",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "stir fry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "shallow fry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "microwave",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "roast",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "smoke",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "sear",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "brine",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "dry",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "ferment",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "marinate",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "pickle",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "season",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "sour",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "sprout",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "cut",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "slice",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "dice",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "grate",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "julienne",
			Description: "measures approximately 1⁄8 by 1⁄8 by 1–2 inches (0.3 cm × 0.3 cm × 3 cm–5 cm)",
			IconPath:    "",
		},
		{
			Name:        "fine julienne",
			Description: "measures approximately 1⁄16 by 1⁄16 by 1–2 inches",
			IconPath:    "",
		},
		{
			Name:        "pont-neuf",
			Description: "measures from 1⁄3 by 1⁄3 by 2 1⁄2–3 inches (1 cm × 1 cm × 6 cm–8 cm) to 3⁄4 by 3⁄4 by 3 inches (2 cm × 2 cm × 8 cm).",
			IconPath:    "",
		},
		{
			Name:        "batonnet",
			Description: "measures approximately 1⁄4 by 1⁄4 by 2–2 1⁄2 inches (0.6 cm × 0.6 cm × 5 cm–6 cm)",
			IconPath:    "",
		},
		{
			Name:        "chiffonade",
			Description: "To roll and cut in sections from 4-10mm in width",
			IconPath:    "",
		},
		{
			Name:        "large dice",
			Description: "sides measuring approximately 3⁄4 inch (20 mm)",
			IconPath:    "",
		},
		{
			Name:        "medium dice",
			Description: "sides measuring approximately 1⁄2 inch",
			IconPath:    "",
		},
		{
			Name:        "small dice",
			Description: "sides measuring approximately 1⁄4 inch (5 mm)",
			IconPath:    "",
		},
		{
			Name:        "brunoise",
			Description: "sides measuring approximately 1⁄8 inch (3 mm)",
			IconPath:    "",
		},
		{
			Name:        "fine brunoise",
			Description: "sides measuring approximately 1⁄16 inch (2 mm)",
			IconPath:    "",
		},
		{
			Name:        "paysanne",
			Description: "1⁄2 by 1⁄2 by 1⁄8 inch (10 mm × 10 mm × 3 mm)",
			IconPath:    "",
		},
		{
			Name:        "lozenge",
			Description: "diamond shape, 1⁄2 by 1⁄2 by 1⁄8 inch (10 mm × 10 mm × 3 mm)",
			IconPath:    "",
		},
		{
			Name:        "tourné",
			Description: "2 inches (50 mm) long with seven faces usually with a bulge in the center portion",
			IconPath:    "",
		},
		{
			Name:        "vacuum seal",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "mince",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "peel",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "shave",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "knead",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "mill",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "grind",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "mix",
			Description: "",
			IconPath:    "",
		},
		{
			Name:        "blend",
			Description: "",
			IconPath:    "",
		},
	}
)
