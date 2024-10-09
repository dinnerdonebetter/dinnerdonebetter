// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductUpdateRequestInput {
  type?: ValidRecipeStepProductType;
  belongsToRecipeStep?: string;
  index?: number;
  isWaste?: boolean;
  quantityNotes?: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions?: string;
  isLiquid?: boolean;
  measurementUnitID?: string;
  name?: string;
  quantity: NumberRange;
  compostable?: boolean;
  containedInVesselIndex?: number;
  storageTemperatureInCelsius: NumberRange;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  type?: ValidRecipeStepProductType;
  belongsToRecipeStep?: string;
  index?: number;
  isWaste?: boolean;
  quantityNotes?: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions?: string;
  isLiquid?: boolean;
  measurementUnitID?: string;
  name?: string;
  quantity: NumberRange;
  compostable?: boolean;
  containedInVesselIndex?: number;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.type = input.type;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.index = input.index;
    this.isWaste = input.isWaste;
    this.quantityNotes = input.quantityNotes;
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions;
    this.isLiquid = input.isLiquid;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.compostable = input.compostable;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
  }
}
