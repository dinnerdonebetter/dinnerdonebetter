// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidVessel } from './ValidVessel';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepVessel {
   archivedAt?: string;
 quantity: NumberRangeWithOptionalMax;
 vesselPreposition: string;
 lastUpdatedAt?: string;
 name: string;
 notes: string;
 recipeStepProductID?: string;
 unavailableAfterStep: boolean;
 belongsToRecipeStep: string;
 createdAt: string;
 id: string;
 vessel?: ValidVessel;

}

export class RecipeStepVessel implements IRecipeStepVessel {
   archivedAt?: string;
 quantity: NumberRangeWithOptionalMax;
 vesselPreposition: string;
 lastUpdatedAt?: string;
 name: string;
 notes: string;
 recipeStepProductID?: string;
 unavailableAfterStep: boolean;
 belongsToRecipeStep: string;
 createdAt: string;
 id: string;
 vessel?: ValidVessel;
constructor(input: Partial<RecipeStepVessel> = {}) {
	 this.archivedAt = input.archivedAt;
 this.quantity = input.quantity = { min: 0 };
 this.vesselPreposition = input.vesselPreposition = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.name = input.name = '';
 this.notes = input.notes = '';
 this.recipeStepProductID = input.recipeStepProductID;
 this.unavailableAfterStep = input.unavailableAfterStep = false;
 this.belongsToRecipeStep = input.belongsToRecipeStep = '';
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.vessel = input.vessel;
}
}