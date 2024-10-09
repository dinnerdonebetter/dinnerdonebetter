// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeUpdateRequestInput {
  eligibleForMeals: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID: string;
  sealOfApproval: boolean;
  slug: string;
  yieldsComponentType: string;
  description: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  source: string;
}

export class RecipeUpdateRequestInput implements IRecipeUpdateRequestInput {
  eligibleForMeals: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID: string;
  sealOfApproval: boolean;
  slug: string;
  yieldsComponentType: string;
  description: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  source: string;
  constructor(input: Partial<RecipeUpdateRequestInput> = {}) {
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.estimatedPortions = input.estimatedPortions || {};
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.sealOfApproval = input.sealOfApproval || false;
    this.slug = input.slug || '';
    this.yieldsComponentType = input.yieldsComponentType || '';
    this.description = input.description || '';
    this.name = input.name || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.portionName = input.portionName || '';
    this.source = input.source || '';
  }
}
