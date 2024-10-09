// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductCreationRequestInput {
  type: ValidRecipeStepProductType;
  isLiquid: boolean;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  name: string;
  quantityNotes: string;
  compostable: boolean;
  containedInVesselIndex: number;
  index: number;
  isWaste: boolean;
  measurementUnitID: string;
}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
  type: ValidRecipeStepProductType;
  isLiquid: boolean;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  name: string;
  quantityNotes: string;
  compostable: boolean;
  containedInVesselIndex: number;
  index: number;
  isWaste: boolean;
  measurementUnitID: string;
  constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
    this.type = input.type || 'ingredient';
    this.isLiquid = input.isLiquid || false;
    this.quantity = input.quantity || { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds || { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.name = input.name || '';
    this.quantityNotes = input.quantityNotes || '';
    this.compostable = input.compostable || false;
    this.containedInVesselIndex = input.containedInVesselIndex || 0;
    this.index = input.index || 0;
    this.isWaste = input.isWaste || false;
    this.measurementUnitID = input.measurementUnitID || '';
  }
}
