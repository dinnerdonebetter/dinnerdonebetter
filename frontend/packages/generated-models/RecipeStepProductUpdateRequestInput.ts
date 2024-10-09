// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidRecipeStepProductType } from './enums';
 import { NumberRange } from './number_range';


export interface IRecipeStepProductUpdateRequestInput {
   isWaste?: boolean;
 storageDurationInSeconds: NumberRange;
 isLiquid?: boolean;
 quantity: NumberRange;
 quantityNotes?: string;
 storageInstructions?: string;
 storageTemperatureInCelsius: NumberRange;
 belongsToRecipeStep?: string;
 containedInVesselIndex?: number;
 measurementUnitID?: string;
 name?: string;
 type?: ValidRecipeStepProductType;
 compostable?: boolean;
 index?: number;

}

export class RecipeStepProductUpdateRequestInput implements IRecipeStepProductUpdateRequestInput {
   isWaste?: boolean;
 storageDurationInSeconds: NumberRange;
 isLiquid?: boolean;
 quantity: NumberRange;
 quantityNotes?: string;
 storageInstructions?: string;
 storageTemperatureInCelsius: NumberRange;
 belongsToRecipeStep?: string;
 containedInVesselIndex?: number;
 measurementUnitID?: string;
 name?: string;
 type?: ValidRecipeStepProductType;
 compostable?: boolean;
 index?: number;
constructor(input: Partial<RecipeStepProductUpdateRequestInput> = {}) {
	 this.isWaste = input.isWaste;
 this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
 this.isLiquid = input.isLiquid;
 this.quantity = input.quantity = { min: 0, max: 0 };
 this.quantityNotes = input.quantityNotes;
 this.storageInstructions = input.storageInstructions;
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.belongsToRecipeStep = input.belongsToRecipeStep;
 this.containedInVesselIndex = input.containedInVesselIndex;
 this.measurementUnitID = input.measurementUnitID;
 this.name = input.name;
 this.type = input.type;
 this.compostable = input.compostable;
 this.index = input.index;
}
}