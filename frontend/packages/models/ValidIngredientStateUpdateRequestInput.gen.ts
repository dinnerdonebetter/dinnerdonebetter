// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums.gen';

export interface IValidIngredientStateUpdateRequestInput {
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
}

export class ValidIngredientStateUpdateRequestInput implements IValidIngredientStateUpdateRequestInput {
  attributeType: ValidIngredientStateAttributeType;
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  constructor(input: Partial<ValidIngredientStateUpdateRequestInput> = {}) {
    this.attributeType = input.attributeType || 'other';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.name = input.name || '';
    this.pastTense = input.pastTense || '';
    this.slug = input.slug || '';
  }
}
