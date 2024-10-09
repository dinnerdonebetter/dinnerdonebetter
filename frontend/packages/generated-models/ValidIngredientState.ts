// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientState {
  iconPath: string;
  id: string;
  name: string;
  slug: string;
  description: string;
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  lastUpdatedAt?: string;
  pastTense: string;
  archivedAt?: string;
}

export class ValidIngredientState implements IValidIngredientState {
  iconPath: string;
  id: string;
  name: string;
  slug: string;
  description: string;
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  lastUpdatedAt?: string;
  pastTense: string;
  archivedAt?: string;
  constructor(input: Partial<ValidIngredientState> = {}) {
    this.iconPath = input.iconPath = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.slug = input.slug = '';
    this.description = input.description = '';
    this.attributeType = input.attributeType = 'other';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.pastTense = input.pastTense = '';
    this.archivedAt = input.archivedAt;
  }
}
