import { RecipeStep } from "./recipeStep";

export class Recipe {
     id: number;
     name: string;
     source: string;
     description: string;
     externalID: string;
     steps: RecipeStep[]
     displayImageURL: string;
     inspiredByRecipeID: number;
     belongsToAccount: number;
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

    constructor(
      id: number = 0,
      name: string = "",
      source: string = "",
      description: string = "",
      externalID: string = "",
      steps: RecipeStep[] = [],
      displayImageURL: string = "",
      inspiredByRecipeID: number = 0,
      belongsToAccount: number = 0,
      createdOn: number = 0,
      lastUpdatedOn?: number,
      archivedOn?: number,
    ) {
         this.id = id;
         this.name = name;
         this.source = source;
         this.description = description;
         this.externalID = externalID;
         this.steps = steps;
         this.displayImageURL = displayImageURL;
         this.inspiredByRecipeID = inspiredByRecipeID;
         this.belongsToAccount = belongsToAccount;
         this.createdOn = createdOn;
         this.lastUpdatedOn = lastUpdatedOn;
         this.archivedOn = archivedOn;
    }
}

export class RecipeList {
     recipes: Recipe[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          recipes: Recipe[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.recipes = recipes;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}