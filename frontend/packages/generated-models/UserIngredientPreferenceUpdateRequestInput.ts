// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserIngredientPreferenceUpdateRequestInput {
  notes?: string;
  rating?: number;
  allergy?: boolean;
  ingredientID?: string;
}

export class UserIngredientPreferenceUpdateRequestInput implements IUserIngredientPreferenceUpdateRequestInput {
  notes?: string;
  rating?: number;
  allergy?: boolean;
  ingredientID?: string;
  constructor(input: Partial<UserIngredientPreferenceUpdateRequestInput> = {}) {
    this.notes = input.notes;
    this.rating = input.rating;
    this.allergy = input.allergy;
    this.ingredientID = input.ingredientID;
  }
}
