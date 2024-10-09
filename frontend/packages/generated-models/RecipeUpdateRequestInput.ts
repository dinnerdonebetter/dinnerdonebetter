// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IRecipeUpdateRequestInput {
   estimatedPortions: OptionalNumberRange;
 inspiredByRecipeID?: string;
 pluralPortionName?: string;
 portionName?: string;
 description?: string;
 eligibleForMeals?: boolean;
 name?: string;
 sealOfApproval?: boolean;
 slug?: string;
 source?: string;
 yieldsComponentType?: string;

}

export class RecipeUpdateRequestInput implements IRecipeUpdateRequestInput {
   estimatedPortions: OptionalNumberRange;
 inspiredByRecipeID?: string;
 pluralPortionName?: string;
 portionName?: string;
 description?: string;
 eligibleForMeals?: boolean;
 name?: string;
 sealOfApproval?: boolean;
 slug?: string;
 source?: string;
 yieldsComponentType?: string;
constructor(input: Partial<RecipeUpdateRequestInput> = {}) {
	 this.estimatedPortions = input.estimatedPortions = {};
 this.inspiredByRecipeID = input.inspiredByRecipeID;
 this.pluralPortionName = input.pluralPortionName;
 this.portionName = input.portionName;
 this.description = input.description;
 this.eligibleForMeals = input.eligibleForMeals;
 this.name = input.name;
 this.sealOfApproval = input.sealOfApproval;
 this.slug = input.slug;
 this.source = input.source;
 this.yieldsComponentType = input.yieldsComponentType;
}
}