// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserIngredientPreferenceUpdateRequestInput {
   allergy: boolean;
 ingredientID: string;
 notes: string;
 rating: number;

}

export class UserIngredientPreferenceUpdateRequestInput implements IUserIngredientPreferenceUpdateRequestInput {
   allergy: boolean;
 ingredientID: string;
 notes: string;
 rating: number;
constructor(input: Partial<UserIngredientPreferenceUpdateRequestInput> = {}) {
	 this.allergy = input.allergy || false;
 this.ingredientID = input.ingredientID || '';
 this.notes = input.notes || '';
 this.rating = input.rating || 0;
}
}