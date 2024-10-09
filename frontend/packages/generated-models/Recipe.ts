// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipeMedia } from './RecipeMedia';
 import { RecipePrepTask } from './RecipePrepTask';
 import { RecipeStep } from './RecipeStep';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipe {
   eligibleForMeals: boolean;
 id: string;
 lastUpdatedAt?: string;
 sealOfApproval: boolean;
 yieldsComponentType: string;
 createdAt: string;
 description: string;
 inspiredByRecipeID?: string;
 media: RecipeMedia;
 name: string;
 slug: string;
 supportingRecipes: Recipe;
 archivedAt?: string;
 createdByUser: string;
 prepTasks: RecipePrepTask;
 source: string;
 estimatedPortions: NumberRangeWithOptionalMax;
 pluralPortionName: string;
 portionName: string;
 steps: RecipeStep;

}

export class Recipe implements IRecipe {
   eligibleForMeals: boolean;
 id: string;
 lastUpdatedAt?: string;
 sealOfApproval: boolean;
 yieldsComponentType: string;
 createdAt: string;
 description: string;
 inspiredByRecipeID?: string;
 media: RecipeMedia;
 name: string;
 slug: string;
 supportingRecipes: Recipe;
 archivedAt?: string;
 createdByUser: string;
 prepTasks: RecipePrepTask;
 source: string;
 estimatedPortions: NumberRangeWithOptionalMax;
 pluralPortionName: string;
 portionName: string;
 steps: RecipeStep;
constructor(input: Partial<Recipe> = {}) {
	 this.eligibleForMeals = input.eligibleForMeals = false;
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.sealOfApproval = input.sealOfApproval = false;
 this.yieldsComponentType = input.yieldsComponentType = '';
 this.createdAt = input.createdAt = '';
 this.description = input.description = '';
 this.inspiredByRecipeID = input.inspiredByRecipeID;
 this.media = input.media = new RecipeMedia();
 this.name = input.name = '';
 this.slug = input.slug = '';
 this.supportingRecipes = input.supportingRecipes = new Recipe();
 this.archivedAt = input.archivedAt;
 this.createdByUser = input.createdByUser = '';
 this.prepTasks = input.prepTasks = new RecipePrepTask();
 this.source = input.source = '';
 this.estimatedPortions = input.estimatedPortions = { min: 0 };
 this.pluralPortionName = input.pluralPortionName = '';
 this.portionName = input.portionName = '';
 this.steps = input.steps = new RecipeStep();
}
}