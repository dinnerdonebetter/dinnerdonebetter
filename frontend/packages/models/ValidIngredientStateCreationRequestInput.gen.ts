// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums.gen';

export interface IValidIngredientStateCreationRequestInput {
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
}

export class ValidIngredientStateCreationRequestInput implements IValidIngredientStateCreationRequestInput {
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  constructor(input: Partial<ValidIngredientStateCreationRequestInput> = {}) {
    this.attributeType = input.attributeType || 'other';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.name = input.name || '';
    this.pastTense = input.pastTense || '';
    this.slug = input.slug || '';
  }
}
