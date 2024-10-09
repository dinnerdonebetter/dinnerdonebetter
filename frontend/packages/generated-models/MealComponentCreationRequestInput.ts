// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponentType } from './enums';

export interface IMealComponentCreationRequestInput {
  recipeID: string;
  recipeScale: number;
  componentType: MealComponentType;
}

export class MealComponentCreationRequestInput implements IMealComponentCreationRequestInput {
  recipeID: string;
  recipeScale: number;
  componentType: MealComponentType;
  constructor(input: Partial<MealComponentCreationRequestInput> = {}) {
    this.recipeID = input.recipeID = '';
    this.recipeScale = input.recipeScale = 0;
    this.componentType = input.componentType = 'unspecified';
  }
}
