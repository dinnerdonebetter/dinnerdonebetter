// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidIngredientGroupUpdateRequestInput {
  description?: string;
  name?: string;
  slug?: string;
}

export class ValidIngredientGroupUpdateRequestInput implements IValidIngredientGroupUpdateRequestInput {
  description?: string;
  name?: string;
  slug?: string;
  constructor(input: Partial<ValidIngredientGroupUpdateRequestInput> = {}) {
    this.description = input.description;
    this.name = input.name;
    this.slug = input.slug;
  }
}
