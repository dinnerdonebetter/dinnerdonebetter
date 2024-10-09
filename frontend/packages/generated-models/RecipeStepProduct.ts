// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProduct {
  quantity: NumberRange;
  storageInstructions: string;
  createdAt: string;
  isWaste: boolean;
  name: string;
  belongsToRecipeStep: string;
  id: string;
  type: ValidRecipeStepProductType;
  lastUpdatedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  archivedAt?: string;
  compostable: boolean;
  isLiquid: boolean;
  containedInVesselIndex?: number;
  index: number;
  quantityNotes: string;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  quantity: NumberRange;
  storageInstructions: string;
  createdAt: string;
  isWaste: boolean;
  name: string;
  belongsToRecipeStep: string;
  id: string;
  type: ValidRecipeStepProductType;
  lastUpdatedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  archivedAt?: string;
  compostable: boolean;
  isLiquid: boolean;
  containedInVesselIndex?: number;
  index: number;
  quantityNotes: string;
  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions = '';
    this.createdAt = input.createdAt = '';
    this.isWaste = input.isWaste = false;
    this.name = input.name = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.id = input.id = '';
    this.type = input.type = 'ingredient';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.measurementUnit = input.measurementUnit;
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.archivedAt = input.archivedAt;
    this.compostable = input.compostable = false;
    this.isLiquid = input.isLiquid = false;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.index = input.index = 0;
    this.quantityNotes = input.quantityNotes = '';
  }
}
