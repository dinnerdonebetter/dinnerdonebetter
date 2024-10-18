// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVessel } from './ValidVessel.gen';
import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipeStepVessel {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vessel: ValidVessel;
  vesselPreposition: string;
}

export class RecipeStepVessel implements IRecipeStepVessel {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID: string;
  unavailableAfterStep: boolean;
  vessel: ValidVessel;
  vesselPreposition: string;
  constructor(input: Partial<RecipeStepVessel> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vessel = input.vessel || new ValidVessel();
    this.vesselPreposition = input.vesselPreposition || '';
  }
}
