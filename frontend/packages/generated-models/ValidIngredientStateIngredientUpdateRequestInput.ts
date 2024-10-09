// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidIngredientStateIngredientUpdateRequestInput {
   notes?: string;
 validIngredientID?: string;
 validIngredientStateID?: string;

}

export class ValidIngredientStateIngredientUpdateRequestInput implements IValidIngredientStateIngredientUpdateRequestInput {
   notes?: string;
 validIngredientID?: string;
 validIngredientStateID?: string;
constructor(input: Partial<ValidIngredientStateIngredientUpdateRequestInput> = {}) {
	 this.notes = input.notes;
 this.validIngredientID = input.validIngredientID;
 this.validIngredientStateID = input.validIngredientStateID;
}
}