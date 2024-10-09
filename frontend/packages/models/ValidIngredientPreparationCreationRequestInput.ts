// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientPreparationCreationRequestInput {
  validIngredientID: string;
  validPreparationID: string;
  notes: string;
}

export class ValidIngredientPreparationCreationRequestInput implements IValidIngredientPreparationCreationRequestInput {
  validIngredientID: string;
  validPreparationID: string;
  notes: string;
  constructor(input: Partial<ValidIngredientPreparationCreationRequestInput> = {}) {
    this.validIngredientID = input.validIngredientID || '';
    this.validPreparationID = input.validPreparationID || '';
    this.notes = input.notes || '';
  }
}
