// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientPreparationCreationRequestInput {
  notes: string;
  validIngredientID: string;
  validPreparationID: string;
}

export class ValidIngredientPreparationCreationRequestInput implements IValidIngredientPreparationCreationRequestInput {
  notes: string;
  validIngredientID: string;
  validPreparationID: string;
  constructor(input: Partial<ValidIngredientPreparationCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.validIngredientID = input.validIngredientID = '';
    this.validPreparationID = input.validPreparationID = '';
  }
}
