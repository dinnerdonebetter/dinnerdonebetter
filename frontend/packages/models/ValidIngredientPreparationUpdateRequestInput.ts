// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientPreparationUpdateRequestInput {
  notes: string;
  validIngredientID: string;
  validPreparationID: string;
}

export class ValidIngredientPreparationUpdateRequestInput implements IValidIngredientPreparationUpdateRequestInput {
  notes: string;
  validIngredientID: string;
  validPreparationID: string;
  constructor(input: Partial<ValidIngredientPreparationUpdateRequestInput> = {}) {
    this.notes = input.notes || '';
    this.validIngredientID = input.validIngredientID || '';
    this.validPreparationID = input.validPreparationID || '';
  }
}
