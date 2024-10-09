// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidInstrument } from './ValidInstrument';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepInstrument {
   quantity: NumberRangeWithOptionalMax;
 id: string;
 instrument?: ValidInstrument;
 lastUpdatedAt?: string;
 name: string;
 optionIndex: number;
 preferenceRank: number;
 recipeStepProductID?: string;
 archivedAt?: string;
 belongsToRecipeStep: string;
 createdAt: string;
 notes: string;
 optional: boolean;

}

export class RecipeStepInstrument implements IRecipeStepInstrument {
   quantity: NumberRangeWithOptionalMax;
 id: string;
 instrument?: ValidInstrument;
 lastUpdatedAt?: string;
 name: string;
 optionIndex: number;
 preferenceRank: number;
 recipeStepProductID?: string;
 archivedAt?: string;
 belongsToRecipeStep: string;
 createdAt: string;
 notes: string;
 optional: boolean;
constructor(input: Partial<RecipeStepInstrument> = {}) {
	 this.quantity = input.quantity = { min: 0 };
 this.id = input.id = '';
 this.instrument = input.instrument;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.name = input.name = '';
 this.optionIndex = input.optionIndex = 0;
 this.preferenceRank = input.preferenceRank = 0;
 this.recipeStepProductID = input.recipeStepProductID;
 this.archivedAt = input.archivedAt;
 this.belongsToRecipeStep = input.belongsToRecipeStep = '';
 this.createdAt = input.createdAt = '';
 this.notes = input.notes = '';
 this.optional = input.optional = false;
}
}