// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserIngredientPreferenceCreationRequestInput {
   allergy: boolean;
 notes: string;
 rating: number;
 validIngredientGroupID: string;
 validIngredientID: string;

}

export class UserIngredientPreferenceCreationRequestInput implements IUserIngredientPreferenceCreationRequestInput {
   allergy: boolean;
 notes: string;
 rating: number;
 validIngredientGroupID: string;
 validIngredientID: string;
constructor(input: Partial<UserIngredientPreferenceCreationRequestInput> = {}) {
	 this.allergy = input.allergy || false;
 this.notes = input.notes || '';
 this.rating = input.rating || 0;
 this.validIngredientGroupID = input.validIngredientGroupID || '';
 this.validIngredientID = input.validIngredientID || '';
}
}