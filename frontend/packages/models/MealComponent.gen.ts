// GENERATED CODE, DO NOT EDIT MANUALLY

import { Recipe } from './Recipe';
import { MealComponentType } from './enums.gen';

export interface IMealComponent {
  componentType: MealComponentType;
  recipe: Recipe;
  recipeScale: number;
}

export class MealComponent implements IMealComponent {
  componentType: MealComponentType;
  recipe: Recipe;
  recipeScale: number;
  constructor(input: Partial<MealComponent> = {}) {
    this.componentType = input.componentType || 'unspecified';
    this.recipe = input.recipe || new Recipe();
    this.recipeScale = input.recipeScale || 0;
  }
}
