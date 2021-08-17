import { RecipeStepIngredient } from "./recipeStepIngredient";

export class RecipeStep {
     id: number;
     temperatureInCelsius?: number;
     externalID: string;
     why: string;
     notes: string;
     index: number;
     prerequisiteStep: number;
     preparationID: number;
     maxEstimatedTimeInSeconds: number;
     minEstimatedTimeInSeconds: number;
     belongsToRecipe: number;
     ingredients: RecipeStepIngredient[];
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

    constructor(
          id: number = 0,
          temperatureInCelsius?: number,
          externalID: string = "",
          why: string = "",
          notes: string = "",
          index: number = 0,
          prerequisiteStep: number = 0,
          preparationID: number = 0,
          maxEstimatedTimeInSeconds: number = 0,
          minEstimatedTimeInSeconds: number = 0,
          belongsToRecipe: number = 0,
          ingredients: RecipeStepIngredient[] = [],
          createdOn: number = 0,
          lastUpdatedOn?: number,
          archivedOn?: number,
    ) {
         this.id = id;
         this.temperatureInCelsius = temperatureInCelsius;
         this.externalID = externalID;
         this.why = why;
         this.notes = notes;
         this.index = index;
         this.prerequisiteStep = prerequisiteStep;
         this.preparationID = preparationID;
         this.maxEstimatedTimeInSeconds = maxEstimatedTimeInSeconds;
         this.minEstimatedTimeInSeconds = minEstimatedTimeInSeconds;
         this.belongsToRecipe = belongsToRecipe;
         this.ingredients = ingredients;
         this.createdOn = createdOn;
         this.lastUpdatedOn = lastUpdatedOn;
         this.archivedOn = archivedOn;
    }
}

export class RecipeStepList {
     recipeSteps: RecipeStep[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          recipeSteps: RecipeStep[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.recipeSteps = recipeSteps;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}