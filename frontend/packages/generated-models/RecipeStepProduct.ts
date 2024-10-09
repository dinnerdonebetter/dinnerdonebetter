// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidMeasurementUnit } from './ValidMeasurementUnit';
 import { ValidRecipeStepProductType } from './enums';
 import { NumberRange } from './number_range';


export interface IRecipeStepProduct {
   isWaste: boolean;
 lastUpdatedAt?: string;
 measurementUnit?: ValidMeasurementUnit;
 storageDurationInSeconds: NumberRange;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 belongsToRecipeStep: string;
 compostable: boolean;
 createdAt: string;
 quantityNotes: string;
 archivedAt?: string;
 isLiquid: boolean;
 name: string;
 quantity: NumberRange;
 containedInVesselIndex?: number;
 id: string;
 index: number;
 type: ValidRecipeStepProductType;

}

export class RecipeStepProduct implements IRecipeStepProduct {
   isWaste: boolean;
 lastUpdatedAt?: string;
 measurementUnit?: ValidMeasurementUnit;
 storageDurationInSeconds: NumberRange;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 belongsToRecipeStep: string;
 compostable: boolean;
 createdAt: string;
 quantityNotes: string;
 archivedAt?: string;
 isLiquid: boolean;
 name: string;
 quantity: NumberRange;
 containedInVesselIndex?: number;
 id: string;
 index: number;
 type: ValidRecipeStepProductType;
constructor(input: Partial<RecipeStepProduct> = {}) {
	 this.isWaste = input.isWaste = false;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.measurementUnit = input.measurementUnit;
 this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
 this.storageInstructions = input.storageInstructions = '';
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.belongsToRecipeStep = input.belongsToRecipeStep = '';
 this.compostable = input.compostable = false;
 this.createdAt = input.createdAt = '';
 this.quantityNotes = input.quantityNotes = '';
 this.archivedAt = input.archivedAt;
 this.isLiquid = input.isLiquid = false;
 this.name = input.name = '';
 this.quantity = input.quantity = { min: 0, max: 0 };
 this.containedInVesselIndex = input.containedInVesselIndex;
 this.id = input.id = '';
 this.index = input.index = 0;
 this.type = input.type = 'ingredient';
}
}