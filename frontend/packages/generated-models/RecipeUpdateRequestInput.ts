// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeUpdateRequestInput {
  name?: string;
  portionName?: string;
  sealOfApproval?: boolean;
  slug?: string;
  source?: string;
  description?: string;
  eligibleForMeals?: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID?: string;
  pluralPortionName?: string;
  yieldsComponentType?: string;
}

export class RecipeUpdateRequestInput implements IRecipeUpdateRequestInput {
  name?: string;
  portionName?: string;
  sealOfApproval?: boolean;
  slug?: string;
  source?: string;
  description?: string;
  eligibleForMeals?: boolean;
  estimatedPortions: OptionalNumberRange;
  inspiredByRecipeID?: string;
  pluralPortionName?: string;
  yieldsComponentType?: string;
  constructor(input: Partial<RecipeUpdateRequestInput> = {}) {
    this.name = input.name;
    this.portionName = input.portionName;
    this.sealOfApproval = input.sealOfApproval;
    this.slug = input.slug;
    this.source = input.source;
    this.description = input.description;
    this.eligibleForMeals = input.eligibleForMeals;
    this.estimatedPortions = input.estimatedPortions = {};
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.pluralPortionName = input.pluralPortionName;
    this.yieldsComponentType = input.yieldsComponentType;
  }
}
