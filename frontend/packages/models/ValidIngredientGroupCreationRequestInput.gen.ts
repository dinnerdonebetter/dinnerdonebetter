// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMemberCreationRequestInput } from './ValidIngredientGroupMemberCreationRequestInput';

export interface IValidIngredientGroupCreationRequestInput {
  description: string;
  members: ValidIngredientGroupMemberCreationRequestInput[];
  name: string;
  slug: string;
}

export class ValidIngredientGroupCreationRequestInput implements IValidIngredientGroupCreationRequestInput {
  description: string;
  members: ValidIngredientGroupMemberCreationRequestInput[];
  name: string;
  slug: string;
  constructor(input: Partial<ValidIngredientGroupCreationRequestInput> = {}) {
    this.description = input.description || '';
    this.members = input.members || [];
    this.name = input.name || '';
    this.slug = input.slug || '';
  }
}
