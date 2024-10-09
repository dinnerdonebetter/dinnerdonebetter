// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredientStateAttributeType } from './enums';

export interface IValidIngredientStateCreationRequestInput {
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
}

export class ValidIngredientStateCreationRequestInput implements IValidIngredientStateCreationRequestInput {
  description: string;
  iconPath: string;
  name: string;
  pastTense: string;
  slug: string;
  attributeType: ValidIngredientStateAttributeType;
  constructor(input: Partial<ValidIngredientStateCreationRequestInput> = {}) {
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.pastTense = input.pastTense = '';
    this.slug = input.slug = '';
    this.attributeType = input.attributeType = 'other';
  }
}
