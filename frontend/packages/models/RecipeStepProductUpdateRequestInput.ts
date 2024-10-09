// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductUpdateRequestInput {
  containedInVesselIndex: number;
  index: number;
  measurementUnitID: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  belongsToRecipeStep: string;
  isWaste: boolean;
  name: string;
  compostable: boolean;
  isLiquid: boolean;
  storageDurationInSeconds: NumberRange;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  containedInVesselIndex: number;
  index: number;
  measurementUnitID: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  belongsToRecipeStep: string;
  isWaste: boolean;
  name: string;
  compostable: boolean;
  isLiquid: boolean;
  storageDurationInSeconds: NumberRange;
  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.containedInVesselIndex = input.containedInVesselIndex || 0;
    this.index = input.index || 0;
    this.measurementUnitID = input.measurementUnitID || '';
    this.quantity = input.quantity || { min: 0, max: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.storageInstructions = input.storageInstructions || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.type = input.type || 'ingredient';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.isWaste = input.isWaste || false;
    this.name = input.name || '';
    this.compostable = input.compostable || false;
    this.isLiquid = input.isLiquid || false;
    this.storageDurationInSeconds = input.storageDurationInSeconds || { min: 0, max: 0 };
  }
}
