// Code generated by gen_typescript. DO NOT EDIT.

import { ValidRecipeStepProductType } from './_unions';
import { NumberRange } from './main';
import { ValidMeasurementUnit } from './validMeasurementUnits';

export interface IRecipeStepProduct {
  createdAt: NonNullable<string>;
  storageTemperatureInCelsius: NonNullable<NumberRange>;
  storageDurationInSeconds: NonNullable<NumberRange>;
  quantity: NonNullable<NumberRange>;
  archivedAt?: string;
  lastUpdatedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  containedInVesselIndex?: number;
  name: NonNullable<string>;
  belongsToRecipeStep: NonNullable<string>;
  type: NonNullable<ValidRecipeStepProductType>;
  id: NonNullable<string>;
  storageInstructions: NonNullable<string>;
  quantityNotes: NonNullable<string>;
  index: NonNullable<number>;
  isWaste: NonNullable<boolean>;
  isLiquid: NonNullable<boolean>;
  compostable: NonNullable<boolean>;
}

export class RecipeStepProduct implements IRecipeStepProduct {
  createdAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  storageTemperatureInCelsius: NonNullable<NumberRange> = { min: 0, max: 0 };
  storageDurationInSeconds: NonNullable<NumberRange> = { min: 0, max: 0 };
  quantity: NonNullable<NumberRange> = { min: 0, max: 0 };
  archivedAt?: string;
  lastUpdatedAt?: string;
  measurementUnit?: ValidMeasurementUnit;
  containedInVesselIndex?: number;
  name: NonNullable<string> = '';
  belongsToRecipeStep: NonNullable<string> = '';
  type: NonNullable<ValidRecipeStepProductType> = 'ingredient';
  id: NonNullable<string> = '';
  storageInstructions: NonNullable<string> = '';
  quantityNotes: NonNullable<string> = '';
  index: NonNullable<number> = 0;
  isWaste: NonNullable<boolean> = false;
  isLiquid: NonNullable<boolean> = false;
  compostable: NonNullable<boolean> = false;

  constructor(input: Partial<RecipeStepProduct> = {}) {
    this.createdAt = input.createdAt ?? '1970-01-01T00:00:00Z';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius ?? { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds ?? { min: 0, max: 0 };
    this.quantity = input.quantity ?? { min: 0, max: 0 };
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.measurementUnit = input.measurementUnit;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.name = input.name ?? '';
    this.belongsToRecipeStep = input.belongsToRecipeStep ?? '';
    this.type = input.type ?? 'ingredient';
    this.id = input.id ?? '';
    this.storageInstructions = input.storageInstructions ?? '';
    this.quantityNotes = input.quantityNotes ?? '';
    this.index = input.index ?? 0;
    this.isWaste = input.isWaste ?? false;
    this.isLiquid = input.isLiquid ?? false;
    this.compostable = input.compostable ?? false;
  }
}

export interface IRecipeStepProductCreationRequestInput {
  storageTemperatureInCelsius: NonNullable<NumberRange>;
  storageDurationInSeconds: NonNullable<NumberRange>;
  quantity: NonNullable<NumberRange>;
  measurementUnitID?: string;
  containedInVesselIndex?: number;
  quantityNotes: NonNullable<string>;
  name: NonNullable<string>;
  storageInstructions: NonNullable<string>;
  type: NonNullable<ValidRecipeStepProductType>;
  index: NonNullable<number>;
  compostable: NonNullable<boolean>;
  isLiquid: NonNullable<boolean>;
  isWaste: NonNullable<boolean>;
}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
  storageTemperatureInCelsius: NonNullable<NumberRange> = { min: 0, max: 0 };
  storageDurationInSeconds: NonNullable<NumberRange> = { min: 0, max: 0 };
  quantity: NonNullable<NumberRange> = { min: 0, max: 0 };
  measurementUnitID?: string;
  containedInVesselIndex?: number;
  quantityNotes: NonNullable<string> = '';
  name: NonNullable<string> = '';
  storageInstructions: NonNullable<string> = '';
  type: NonNullable<ValidRecipeStepProductType> = 'ingredient';
  index: NonNullable<number> = 0;
  compostable: NonNullable<boolean> = false;
  isLiquid: NonNullable<boolean> = false;
  isWaste: NonNullable<boolean> = false;

  constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius ?? { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds ?? { min: 0, max: 0 };
    this.quantity = input.quantity ?? { min: 0, max: 0 };
    this.measurementUnitID = input.measurementUnitID;
    this.containedInVesselIndex = input.containedInVesselIndex;
    this.quantityNotes = input.quantityNotes ?? '';
    this.name = input.name ?? '';
    this.storageInstructions = input.storageInstructions ?? '';
    this.type = input.type ?? 'ingredient';
    this.index = input.index ?? 0;
    this.compostable = input.compostable ?? false;
    this.isLiquid = input.isLiquid ?? false;
    this.isWaste = input.isWaste ?? false;
  }
}

export interface IRecipeStepProductUpdateRequestInput {
  name?: string;
  type?: ValidRecipeStepProductType;
  measurementUnitID?: string;
  quantityNotes?: string;
  belongsToRecipeStep?: string;
  storageTemperatureInCelsius: NonNullable<NumberRange>;
  storageDurationInSeconds: NonNullable<NumberRange>;
  quantity: NonNullable<NumberRange>;
  compostable?: boolean;
  storageInstructions?: string;
  isLiquid?: boolean;
  isWaste?: boolean;
  index?: number;
  containedInVesselIndex?: number;
}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
  name?: string;
  type?: ValidRecipeStepProductType = 'ingredient';
  measurementUnitID?: string;
  quantityNotes?: string;
  belongsToRecipeStep?: string;
  storageTemperatureInCelsius: NonNullable<NumberRange> = { min: 0, max: 0 };
  storageDurationInSeconds: NonNullable<NumberRange> = { min: 0, max: 0 };
  quantity: NonNullable<NumberRange> = { min: 0, max: 0 };
  compostable?: boolean = false;
  storageInstructions?: string;
  isLiquid?: boolean = false;
  isWaste?: boolean = false;
  index?: number;
  containedInVesselIndex?: number;

  constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
    this.name = input.name;
    this.type = input.type ?? 'ingredient';
    this.measurementUnitID = input.measurementUnitID;
    this.quantityNotes = input.quantityNotes;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius ?? { min: 0, max: 0 };
    this.storageDurationInSeconds = input.storageDurationInSeconds ?? { min: 0, max: 0 };
    this.quantity = input.quantity ?? { min: 0, max: 0 };
    this.compostable = input.compostable ?? false;
    this.storageInstructions = input.storageInstructions;
    this.isLiquid = input.isLiquid ?? false;
    this.isWaste = input.isWaste ?? false;
    this.index = input.index;
    this.containedInVesselIndex = input.containedInVesselIndex;
  }
}
