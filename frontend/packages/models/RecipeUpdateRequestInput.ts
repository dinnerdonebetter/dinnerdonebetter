// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeUpdateRequestInput {
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  sealOfApproval: boolean;
  slug: string;
  source: string;
  yieldsComponentType: string;
}

export class RecipeUpdateRequestInput implements IRecipeUpdateRequestInput {
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  sealOfApproval: boolean;
  slug: string;
  source: string;
  yieldsComponentType: string;
  constructor(input: Partial<RecipeUpdateRequestInput> = {}) {
    this.description = input.description || '';
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.estimatedPortions = input.estimatedPortions || {};
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.name = input.name || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.portionName = input.portionName || '';
    this.sealOfApproval = input.sealOfApproval || false;
    this.slug = input.slug || '';
    this.source = input.source || '';
    this.yieldsComponentType = input.yieldsComponentType || '';
  }
}
