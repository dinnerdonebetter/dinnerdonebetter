// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserIngredientPreferenceCreationRequestInput {
  validIngredientID: string;
  allergy: boolean;
  notes: string;
  rating: number;
  validIngredientGroupID: string;
}

export class UserIngredientPreferenceCreationRequestInput implements IUserIngredientPreferenceCreationRequestInput {
  validIngredientID: string;
  allergy: boolean;
  notes: string;
  rating: number;
  validIngredientGroupID: string;
  constructor(input: Partial<UserIngredientPreferenceCreationRequestInput> = {}) {
    this.validIngredientID = input.validIngredientID = '';
    this.allergy = input.allergy = false;
    this.notes = input.notes = '';
    this.rating = input.rating = 0;
    this.validIngredientGroupID = input.validIngredientGroupID = '';
  }
}
