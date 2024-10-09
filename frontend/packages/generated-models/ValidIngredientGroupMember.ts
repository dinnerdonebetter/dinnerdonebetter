// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IValidIngredientGroupMember {
  archivedAt?: string;
  belongsToGroup: string;
  createdAt: string;
  id: string;
  validIngredient: ValidIngredient;
}

export class ValidIngredientGroupMember implements IValidIngredientGroupMember {
  archivedAt?: string;
  belongsToGroup: string;
  createdAt: string;
  id: string;
  validIngredient: ValidIngredient;
  constructor(input: Partial<ValidIngredientGroupMember> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToGroup = input.belongsToGroup = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.validIngredient = input.validIngredient = new ValidIngredient();
  }
}
