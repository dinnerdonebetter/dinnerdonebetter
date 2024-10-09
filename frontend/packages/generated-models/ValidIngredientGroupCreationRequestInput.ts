// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMemberCreationRequestInput } from './ValidIngredientGroupMemberCreationRequestInput';

export interface IValidIngredientGroupCreationRequestInput {
  slug: string;
  description: string;
  members: ValidIngredientGroupMemberCreationRequestInput;
  name: string;
}

export class ValidIngredientGroupCreationRequestInput implements IValidIngredientGroupCreationRequestInput {
  slug: string;
  description: string;
  members: ValidIngredientGroupMemberCreationRequestInput;
  name: string;
  constructor(input: Partial<ValidIngredientGroupCreationRequestInput> = {}) {
    this.slug = input.slug = '';
    this.description = input.description = '';
    this.members = input.members = new ValidIngredientGroupMemberCreationRequestInput();
    this.name = input.name = '';
  }
}
