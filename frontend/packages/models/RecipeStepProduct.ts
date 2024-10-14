// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProduct {
  archivedAt: string;
  belongsToRecipeStep: string;
  compostable: boolean;
  containedInVesselIndex: number;
  createdAt: string;
  id: string;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  archivedAt: string;
  belongsToRecipeStep: string;
  compostable: boolean;
  containedInVesselIndex: number;
  createdAt: string;
  id: string;
  index: number;
  isLiquid: boolean;
  isWaste: boolean;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  quantity: NumberRange;
  quantityNotes: string;
  storageDurationInSeconds: NumberRange;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.compostable = input.compostable || false;
    this.containedInVesselIndex = input.containedInVesselIndex || 0;
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.index = input.index || 0;
    this.isLiquid = input.isLiquid || false;
    this.isWaste = input.isWaste || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.name = input.name || '';
    this.quantity = input.quantity || { min: 0, max: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.storageDurationInSeconds = input.storageDurationInSeconds || { min: 0, max: 0 };
    this.storageInstructions = input.storageInstructions || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.type = input.type || 'ingredient';
  }
}
