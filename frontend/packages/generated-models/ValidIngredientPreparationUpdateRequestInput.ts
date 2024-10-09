// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientPreparationUpdateRequestInput {
  validPreparationID?: string;
  notes?: string;
  validIngredientID?: string;
}

export class ValidIngredientPreparationUpdateRequestInput implements IValidIngredientPreparationUpdateRequestInput {
  validPreparationID?: string;
  notes?: string;
  validIngredientID?: string;
  constructor(input: Partial<ValidIngredientPreparationUpdateRequestInput> = {}) {
    this.validPreparationID = input.validPreparationID;
    this.notes = input.notes;
    this.validIngredientID = input.validIngredientID;
  }
}
