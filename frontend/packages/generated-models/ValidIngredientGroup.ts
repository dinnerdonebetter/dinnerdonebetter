// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMember } from './ValidIngredientGroupMember';

export interface IValidIngredientGroup {
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
  name: string;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
}

export class ValidIngredientGroup implements IValidIngredientGroup {
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
  name: string;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  constructor(input: Partial<ValidIngredientGroup> = {}) {
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.members = input.members = new ValidIngredientGroupMember();
    this.name = input.name = '';
    this.slug = input.slug = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
  }
}
