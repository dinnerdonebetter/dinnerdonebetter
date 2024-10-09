// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProduct {
  id: string;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnit?: ValidMeasurementUnit;
  quantity: NumberRange;
  archivedAt?: string;
  createdAt: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  containedInVesselIndex?: number;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  belongsToRecipeStep: string;
  compostable: boolean;
  lastUpdatedAt?: string;
  name: string;
  quantityNotes: string;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  id: string;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  measurementUnit?: ValidMeasurementUnit;
  quantity: NumberRange;
  archivedAt?: string;
  createdAt: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  containedInVesselIndex?: number;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  belongsToRecipeStep: string;
  compostable: boolean;
  lastUpdatedAt?: string;
  name: string;
  quantityNotes: string;
  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.id = input.id = '';
    this.index = input.index = 0;
    this.isLiquid = input.isLiquid = false;
    this.isWaste = input.isWaste = false;
    this.measurementUnit = input.measurementUnit;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions = '';
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.type = input.type = 'ingredient';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.compostable = input.compostable = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.quantityNotes = input.quantityNotes = '';
  }
}
