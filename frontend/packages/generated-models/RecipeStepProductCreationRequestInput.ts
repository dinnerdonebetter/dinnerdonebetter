// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidRecipeStepProductType } from './enums';
 import { NumberRange } from './number_range';


export interface IRecipeStepProductCreationRequestInput {
   containedInVesselIndex?: number;
 isLiquid: boolean;
 isWaste: boolean;
 name: string;
 storageDurationInSeconds: NumberRange;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 type: ValidRecipeStepProductType;
 compostable: boolean;
 index: number;
 measurementUnitID?: string;
 quantity: NumberRange;
 quantityNotes: string;

}

export class RecipeStepProductCreationRequestInput implements IRecipeStepProductCreationRequestInput {
   containedInVesselIndex?: number;
 isLiquid: boolean;
 isWaste: boolean;
 name: string;
 storageDurationInSeconds: NumberRange;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 type: ValidRecipeStepProductType;
 compostable: boolean;
 index: number;
 measurementUnitID?: string;
 quantity: NumberRange;
 quantityNotes: string;
constructor(input: Partial<RecipeStepProductCreationRequestInput> = {}) {
	 this.containedInVesselIndex = input.containedInVesselIndex;
 this.isLiquid = input.isLiquid = false;
 this.isWaste = input.isWaste = false;
 this.name = input.name = '';
 this.storageDurationInSeconds = input.storageDurationInSeconds = { min: 0, max: 0 };
 this.storageInstructions = input.storageInstructions = '';
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.type = input.type = 'ingredient';
 this.compostable = input.compostable = false;
 this.index = input.index = 0;
 this.measurementUnitID = input.measurementUnitID;
 this.quantity = input.quantity = { min: 0, max: 0 };
 this.quantityNotes = input.quantityNotes = '';
}
}