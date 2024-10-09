// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientGroupMember } from './ValidIngredientGroupMember';

export interface IValidIngredientGroup {
  description: string;
  id: string;
  lastUpdatedAt: string;
  members: ValidIngredientGroupMember[];
  name: string;
  slug: string;
  archivedAt: string;
  createdAt: string;
}

export class ValidIngredientGroup implements IValidIngredientGroup {
  description: string;
  id: string;
  lastUpdatedAt: string;
  members: ValidIngredientGroupMember[];
  name: string;
  slug: string;
  archivedAt: string;
  createdAt: string;
  constructor(input: Partial<ValidIngredientGroup> = {}) {
    this.description = input.description || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.members = input.members || [];
    this.name = input.name || '';
    this.slug = input.slug || '';
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
  }
}
