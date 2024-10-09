// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientState {
  archivedAt: string;
  createdAt: string;
  description: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
  iconPath: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
}

export class ValidIngredientState implements IValidIngredientState {
  archivedAt: string;
  createdAt: string;
  description: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
  iconPath: string;
  id: string;
  lastUpdatedAt: string;
  name: string;
  constructor(input: Partial<ValidIngredientState> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.pastTense = input.pastTense || '';
    this.slug = input.slug || '';
    this.attributeType = input.attributeType || 'other';
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
  }
}
