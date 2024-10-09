// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVessel } from './ValidVessel';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepVessel {
  lastUpdatedAt?: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  createdAt: string;
  belongsToRecipeStep: string;
  id: string;
  name: string;
  recipeStepProductID?: string;
  vessel?: ValidVessel;
  vesselPreposition: string;
  archivedAt?: string;
}

export class RecipeStepVessel implements IRecipeStepVessel {
  lastUpdatedAt?: string;
  notes: string;
  quantity: NumberRangeWithOptionalMax;
  unavailableAfterStep: boolean;
  createdAt: string;
  belongsToRecipeStep: string;
  id: string;
  name: string;
  recipeStepProductID?: string;
  vessel?: ValidVessel;
  vesselPreposition: string;
  archivedAt?: string;
  constructor(input: Partial<RecipeStepVessel> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.quantity = input.quantity = { min: 0 };
    this.unavailableAfterStep = input.unavailableAfterStep = false;
    this.createdAt = input.createdAt = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.recipeStepProductID = input.recipeStepProductID;
    this.vessel = input.vessel;
    this.vesselPreposition = input.vesselPreposition = '';
    this.archivedAt = input.archivedAt;
  }
}
