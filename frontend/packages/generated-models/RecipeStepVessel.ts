// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVessel } from './ValidVessel';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVessel {
  lastUpdatedAt?: string;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  belongsToRecipeStep: string;
  createdAt: string;
  notes: string;
  recipeStepProductID?: string;
  vessel?: ValidVessel;
  archivedAt?: string;
  id: string;
}

export class RecipeStepVessel implements IRecipeStepVessel {
  lastUpdatedAt?: string;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  vesselPreposition: string;
  belongsToRecipeStep: string;
  createdAt: string;
  notes: string;
  recipeStepProductID?: string;
  vessel?: ValidVessel;
  archivedAt?: string;
  id: string;
  constructor(input: Partial<RecipeStepVessel> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.quantity = input.quantity = { min: 0 };
    this.unavailableAfterStep = input.unavailableAfterStep = false;
    this.vesselPreposition = input.vesselPreposition = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.createdAt = input.createdAt = '';
    this.notes = input.notes = '';
    this.recipeStepProductID = input.recipeStepProductID;
    this.vessel = input.vessel;
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
  }
}
