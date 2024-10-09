// GENERATED CODE, DO NOT EDIT MANUALLY

import { Recipe } from './Recipe';
import { MealComponentType } from './enums';

export interface IMealComponent {
  recipe: Recipe;
  recipeScale: number;
  componentType: MealComponentType;
}

export class MealComponent implements IMealComponent {
  recipe: Recipe;
  recipeScale: number;
  componentType: MealComponentType;
  constructor(input: Partial<MealComponent> = {}) {
    this.recipe = input.recipe = new Recipe();
    this.recipeScale = input.recipeScale = 0;
    this.componentType = input.componentType = 'unspecified';
  }
}
