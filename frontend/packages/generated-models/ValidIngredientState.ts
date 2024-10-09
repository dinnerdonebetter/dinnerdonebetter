// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientState {
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  archivedAt?: string;
  id: string;
  lastUpdatedAt?: string;
}

export class ValidIngredientState implements IValidIngredientState {
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  archivedAt?: string;
  id: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<ValidIngredientState> = {}) {
    this.attributeType = input.attributeType = 'other';
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.pastTense = input.pastTense = '';
    this.slug = input.slug = '';
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
