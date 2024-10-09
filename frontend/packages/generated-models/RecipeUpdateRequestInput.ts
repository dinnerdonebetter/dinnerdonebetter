// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeUpdateRequestInput {
  description?: string;
  estimatedPortions: OptionalNumberRange;
  slug?: string;
  source?: string;
  yieldsComponentType?: string;
  eligibleForMeals?: boolean;
  inspiredByRecipeID?: string;
  name?: string;
  pluralPortionName?: string;
  portionName?: string;
  sealOfApproval?: boolean;
}

export class RecipeUpdateRequestInput implements IRecipeUpdateRequestInput {
  description?: string;
  estimatedPortions: OptionalNumberRange;
  slug?: string;
  source?: string;
  yieldsComponentType?: string;
  eligibleForMeals?: boolean;
  inspiredByRecipeID?: string;
  name?: string;
  pluralPortionName?: string;
  portionName?: string;
  sealOfApproval?: boolean;
  constructor(input: Partial<RecipeUpdateRequestInput> = {}) {
    this.description = input.description;
    this.estimatedPortions = input.estimatedPortions = {};
    this.slug = input.slug;
    this.source = input.source;
    this.yieldsComponentType = input.yieldsComponentType;
    this.eligibleForMeals = input.eligibleForMeals;
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.name = input.name;
    this.pluralPortionName = input.pluralPortionName;
    this.portionName = input.portionName;
    this.sealOfApproval = input.sealOfApproval;
  }
}
