// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums.gen';
import { NumberRange } from './number_range.gen';

export interface IRecipeStepProductUpdateRequestInput {
  belongsToRecipeStep: string;
  compostable: boolean;
  containedInVesselIndex: number;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnitID: string;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  belongsToRecipeStep: string;
  compostable: boolean;
  containedInVesselIndex: number;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnitID: string;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.compostable = input.compostable || false;
    this.containedInVesselIndex = input.containedInVesselIndex || 0;
    this.index = input.index || 0;
    this.isLiquid = input.isLiquid || false;
    this.isWaste = input.isWaste || false;
    this.measurementUnitID = input.measurementUnitID || '';
    this.name = input.name || '';
    this.quantity = input.quantity || { min: 0, max: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.storageDurationInSeconds = input.storageDurationInSeconds || { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.type = input.type || 'ingredient';
  }
}
