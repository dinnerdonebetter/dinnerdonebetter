// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMember } from './ValidIngredientGroupMember';

export interface IValidIngredientGroup {
  name: string;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
}

export class ValidIngredientGroup implements IValidIngredientGroup {
  name: string;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
  constructor(input: Partial<ValidIngredientGroup> = {}) {
    this.name = input.name = '';
    this.slug = input.slug = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.members = input.members = new ValidIngredientGroupMember();
  }
}
