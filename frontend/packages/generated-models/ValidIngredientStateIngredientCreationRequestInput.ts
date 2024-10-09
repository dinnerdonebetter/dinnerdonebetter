// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidIngredientStateIngredientCreationRequestInput {
   notes: string;
 validIngredientID: string;
 validIngredientStateID: string;

}

export class ValidIngredientStateIngredientCreationRequestInput implements IValidIngredientStateIngredientCreationRequestInput {
   notes: string;
 validIngredientID: string;
 validIngredientStateID: string;
constructor(input: Partial<ValidIngredientStateIngredientCreationRequestInput> = {}) {
	 this.notes = input.notes = '';
 this.validIngredientID = input.validIngredientID = '';
 this.validIngredientStateID = input.validIngredientStateID = '';
}
}