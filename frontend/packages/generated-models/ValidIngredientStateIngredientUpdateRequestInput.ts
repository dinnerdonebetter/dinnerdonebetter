// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientStateIngredientUpdateRequestInput {
  validIngredientStateID?: string;
  notes?: string;
  validIngredientID?: string;
}

export class ValidIngredientStateIngredientUpdateRequestInput
  implements IValidIngredientStateIngredientUpdateRequestInput
{
  validIngredientStateID?: string;
  notes?: string;
  validIngredientID?: string;
  constructor(input: Partial<ValidIngredientStateIngredientUpdateRequestInput> = {}) {
    this.validIngredientStateID = input.validIngredientStateID;
    this.notes = input.notes;
    this.validIngredientID = input.validIngredientID;
  }
}
