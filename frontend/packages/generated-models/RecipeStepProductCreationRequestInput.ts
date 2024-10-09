// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductCreationRequestInput {
  quantity: NumberRange;
  storageInstructions: string;
  type: ValidRecipeStepProductType;
  compostable: boolean;
  index: number;
  isLiquid: boolean;
  measurementUnitID?: string;
  name: string;
  containedInVesselIndex?: number;
  isWaste: boolean;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
  quantity: NumberRange;
  storageInstructions: string;
  type: ValidRecipeStepProductType;
  compostable: boolean;
  index: number;
  isLiquid: boolean;
  measurementUnitID?: string;
  name: string;
  containedInVesselIndex?: number;
  isWaste: boolean;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions = '';
    this.type = input.type = 'ingredient';
    this.compostable = input.compostable = false;
    this.index = input.index = 0;
    this.isLiquid = input.isLiquid = false;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name = '';
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.isWaste = input.isWaste = false;
    this.quantityNotes = input.quantityNotes = '';
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
  }
}
