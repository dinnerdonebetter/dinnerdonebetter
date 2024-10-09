// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientState {
  description: string;
  iconPath: string;
  lastUpdatedAt?: string;
  name: string;
  pastTense: string;
  slug: string;
  archivedAt?: string;
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  id: string;
}

export class ValidIngredientState implements IValidIngredientState {
  description: string;
  iconPath: string;
  lastUpdatedAt?: string;
  name: string;
  pastTense: string;
  slug: string;
  archivedAt?: string;
  attributeType: ValidIngredientStateAttributeType;
  createdAt: string;
  id: string;
  constructor(input: Partial<ValidIngredientState> = {}) {
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.pastTense = input.pastTense = '';
    this.slug = input.slug = '';
    this.archivedAt = input.archivedAt;
    this.attributeType = input.attributeType = 'other';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
  }
}
