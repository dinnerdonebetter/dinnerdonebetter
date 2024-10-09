// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidRecipeStepProductType } from './enums';
import { NumberRange } from './number_range';

export interface IRecipeStepProduct {
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  index: number;
  isLiquid: boolean;
  lastUpdatedAt: string;
  quantity: NumberRange;
  quantityNotes: string;
  belongsToRecipeStep: string;
  createdAt: string;
  measurementUnit: ValidMeasurementUnit;
  storageDurationInSeconds: NumberRange;
  containedInVesselIndex: number;
  id: string;
  isWaste: boolean;
  archivedAt: string;
  compostable: boolean;
  name: string;
  storageInstructions: string;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  storageTemperatureInCelsius: NumberRange;
  type: ValidRecipeStepProductType;
  index: number;
  isLiquid: boolean;
  lastUpdatedAt: string;
  quantity: NumberRange;
  quantityNotes: string;
  belongsToRecipeStep: string;
  createdAt: string;
  measurementUnit: ValidMeasurementUnit;
  storageDurationInSeconds: NumberRange;
  containedInVesselIndex: number;
  id: string;
  isWaste: boolean;
  archivedAt: string;
  compostable: boolean;
  name: string;
  storageInstructions: string;
  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.type = input.type || 'ingredient';
    this.index = input.index || 0;
    this.isLiquid = input.isLiquid || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.quantity = input.quantity || { min: 0, max: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.storageDurationInSeconds = input.storageDurationInSeconds || { min: 0, max: 0 };
    this.containedInVesselIndex = input.containedInVesselIndex || 0;
    this.id = input.id || '';
    this.isWaste = input.isWaste || false;
    this.archivedAt = input.archivedAt || '';
    this.compostable = input.compostable || false;
    this.name = input.name || '';
    this.storageInstructions = input.storageInstructions || '';
  }
}
