// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProductUpdateRequestInput {
  quantityNotes?: string;
  isWaste?: boolean;
  index?: number;
  belongsToRecipeStep?: string;
  isLiquid?: boolean;
  containedInVesselIndex?: number;
  measurementUnitID?: string;
  name?: string;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions?: string;
  storageTemperatureInCelsius: NumberRange;
  compostable?: boolean;
  type?: ValidRecipeStepProductType;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  quantityNotes?: string;
  isWaste?: boolean;
  index?: number;
  belongsToRecipeStep?: string;
  isLiquid?: boolean;
  containedInVesselIndex?: number;
  measurementUnitID?: string;
  name?: string;
  quantity: NumberRange;
  storageDurationInSeconds: NumberRange;
  storageInstructions?: string;
  storageTemperatureInCelsius: NumberRange;
  compostable?: boolean;
  type?: ValidRecipeStepProductType;
  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.quantityNotes = input.quantityNotes;
    this.isWaste = input.isWaste;
    this.index = input.index;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.isLiquid = input.isLiquid;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.compostable = input.compostable;
    this.type = input.type;
  }
}
