// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVessel } from './ValidVessel';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVessel {
  vesselPreposition: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  notes: string;
  vessel?: ValidVessel;
  lastUpdatedAt?: string;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  unavailableAfterStep: boolean;
}

export class RecipeStepVessel implements IRecipeStepVessel {
  vesselPreposition: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  notes: string;
  vessel?: ValidVessel;
  lastUpdatedAt?: string;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  recipeStepProductID?: string;
  unavailableAfterStep: boolean;
  constructor(input: Partial<RecipeStepVessel> = {}) {
    this.vesselPreposition = input.vesselPreposition = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.vessel = input.vessel;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.quantity = input.quantity = { min: 0 };
    this.recipeStepProductID = input.recipeStepProductID;
    this.unavailableAfterStep = input.unavailableAfterStep = false;
  }
}
