// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserIngredientPreferenceCreationRequestInput {
   notes: string;
 rating: number;
 validIngredientGroupID: string;
 validIngredientID: string;
 allergy: boolean;

}

export class UserIngredientPreferenceCreationRequestInput implements IUserIngredientPreferenceCreationRequestInput {
   notes: string;
 rating: number;
 validIngredientGroupID: string;
 validIngredientID: string;
 allergy: boolean;
constructor(input: Partial<UserIngredientPreferenceCreationRequestInput> = {}) {
	 this.notes = input.notes = '';
 this.rating = input.rating = 0;
 this.validIngredientGroupID = input.validIngredientGroupID = '';
 this.validIngredientID = input.validIngredientID = '';
 this.allergy = input.allergy = false;
}
}