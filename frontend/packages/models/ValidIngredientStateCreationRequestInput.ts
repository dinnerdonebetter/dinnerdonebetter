// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientStateCreationRequestInput {
  name: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
}

export class ValidIngredientStateCreationRequestInput implements IValidIngredientStateCreationRequestInput {
  name: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
  constructor(input: Partial<ValidIngredientStateCreationRequestInput> = {}) {
    this.name = input.name || '';
    this.pastTense = input.pastTense || '';
    this.slug = input.slug || '';
    this.attributeType = input.attributeType || 'other';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
  }
}
