// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductCreationRequestInput {
  isLiquid: boolean;
  measurementUnitID?: string;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  compostable: boolean;
  containedInVesselIndex?: number;
  index: number;
  isWaste: boolean;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
  isLiquid: boolean;
  measurementUnitID?: string;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  compostable: boolean;
  containedInVesselIndex?: number;
  index: number;
  isWaste: boolean;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
    this.isLiquid = input.isLiquid = false;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name = '';
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.quantityNotes = input.quantityNotes = '';
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.compostable = input.compostable = false;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.index = input.index = 0;
    this.isWaste = input.isWaste = false;
    this.storageInstructions = input.storageInstructions = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.type = input.type = 'ingredient';
  }
}
