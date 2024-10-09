// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductCreationRequestInput {
  containedInVesselIndex?: number;
  index: number;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  compostable: boolean;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnitID?: string;
  name: string;
  quantityNotes: string;
  type: ValidRecipeStepProductType;
}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
  containedInVesselIndex?: number;
  index: number;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  compostable: boolean;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnitID?: string;
  name: string;
  quantityNotes: string;
  type: ValidRecipeStepProductType;
  constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.index = input.index = 0;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.compostable = input.compostable = false;
    this.isLiquid = input.isLiquid = false;
    this.isWaste = input.isWaste = false;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name = '';
    this.quantityNotes = input.quantityNotes = '';
    this.type = input.type = 'ingredient';
  }
}
