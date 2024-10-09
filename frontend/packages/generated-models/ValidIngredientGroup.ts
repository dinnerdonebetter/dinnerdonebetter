// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMember } from './ValidIngredientGroupMember';

export interface IValidIngredientGroup {
  archivedAt?: string;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
  name: string;
  slug: string;
}

export class ValidIngredientGroup implements IValidIngredientGroup {
  archivedAt?: string;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  members: ValidIngredientGroupMember;
  name: string;
  slug: string;
  constructor(input: Partial<ValidIngredientGroup> = {}) {
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.members = input.members = new ValidIngredientGroupMember();
    this.name = input.name = '';
    this.slug = input.slug = '';
  }
}
