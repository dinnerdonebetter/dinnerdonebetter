// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProduct {
  archivedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  isWaste: boolean;
  name: string;
  quantityNotes: string;
  storageInstructions: string;
  type: ValidRecipeStepProductType;
  compostable: boolean;
  createdAt: string;
  id: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  belongsToRecipeStep: string;
  isLiquid: boolean;
  quantity: NumberRange;
  containedInVesselIndex?: number;
  index: number;
  lastUpdatedAt?: string;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  archivedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  isWaste: boolean;
  name: string;
  quantityNotes: string;
  storageInstructions: string;
  type: ValidRecipeStepProductType;
  compostable: boolean;
  createdAt: string;
  id: string;
  storageDurationInSeconds: NumberRange;
  storageTemperatureInCelsius: NumberRange;
  belongsToRecipeStep: string;
  isLiquid: boolean;
  quantity: NumberRange;
  containedInVesselIndex?: number;
  index: number;
  lastUpdatedAt?: string;
  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.archivedAt = input.archivedAt;
    this.measurementUnit = input.measurementUnit;
    this.isWaste = input.isWaste = false;
    this.name = input.name = '';
    this.quantityNotes = input.quantityNotes = '';
    this.storageInstructions = input.storageInstructions = '';
    this.type = input.type = 'ingredient';
    this.compostable = input.compostable = false;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.isLiquid = input.isLiquid = false;
    this.quantity = input.quantity = { min: 0, max: 0 };
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.index = input.index = 0;
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
