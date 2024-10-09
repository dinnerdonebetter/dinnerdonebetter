// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVessel } from './ValidVessel';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVessel {
  archivedAt: string;
  belongsToRecipeStep: string;
  id: string;
  recipeStepProductID: string;
  vesselPreposition: string;
  createdAt: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vessel: ValidVessel;
}

export class RecipeStepVessel implements IRecipeStepVessel {
  archivedAt: string;
  belongsToRecipeStep: string;
  id: string;
  recipeStepProductID: string;
  vesselPreposition: string;
  createdAt: string;
  lastUpdatedAt: string;
  name: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vessel: ValidVessel;
  constructor(input: Partial<RecipeStepVessel> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.id = input.id || '';
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.vesselPreposition = input.vesselPreposition || '';
    this.createdAt = input.createdAt || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || { min: 0 };
    this.unavailableAfterStep = input.unavailableAfterStep || false;
    this.vessel = input.vessel || new ValidVessel();
  }
}
