// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductUpdateRequestInput {
  compostable?: boolean;
  name?: string;
  quantityNotes?: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  type?: ValidRecipeStepProductType;
  containedInVesselIndex?: number;
  index?: number;
  measurementUnitID?: string;
  quantity: NumberRange;
  storageInstructions?: string;
  belongsToRecipeStep?: string;
  isLiquid?: boolean;
  isWaste?: boolean;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  compostable?: boolean;
  name?: string;
  quantityNotes?: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  type?: ValidRecipeStepProductType;
  containedInVesselIndex?: number;
  index?: number;
  measurementUnitID?: string;
  quantity: NumberRange;
  storageInstructions?: string;
  belongsToRecipeStep?: string;
  isLiquid?: boolean;
  isWaste?: boolean;
  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.compostable = input.compostable;
    this.name = input.name;
    this.quantityNotes = input.quantityNotes;
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.type = input.type;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.index = input.index;
    this.measurementUnitID = input.measurementUnitID;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.isLiquid = input.isLiquid;
    this.isWaste = input.isWaste;
  }
}
