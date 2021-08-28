import { ValidIngredient } from "./validIngredient";

export type QuantityType = 'grams' | 'fl. oz';

export class RecipeStepIngredient {
     id: number;
     name: string;
     ingredientID?: number;
     ingredientNotes: string;
     quantityType: QuantityType;
     quantityNotes: string;
     quantityValue: number;
     productOfRecipeStep: boolean;
     belongsToRecipeStep?: number;
     externalID: string;
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

     constructor(
         id: number = 0,
         name: string = '',
         ingredientID?: number,
         ingredientNotes: string = '',
         quantityType: QuantityType = 'grams',
         quantityNotes: string = '',
         quantityValue: number = 0,
         productOfRecipeStep: boolean = false,
         belongsToRecipeStep?: number,
         externalID: string = '',
         createdOn: number = 0,
         lastUpdatedOn?: number,
         archivedOn?: number,
     ) {
          this.id = id;
          this.name = name;
          this.ingredientID = ingredientID;
          this.ingredientNotes = ingredientNotes;
          this.quantityType = quantityType;
          this.quantityNotes = quantityNotes;
          this.quantityValue = quantityValue;
          this.productOfRecipeStep = productOfRecipeStep;
          this.belongsToRecipeStep = belongsToRecipeStep;
          this.externalID = externalID;
          this.createdOn = createdOn;
          this.lastUpdatedOn = lastUpdatedOn;
          this.archivedOn = archivedOn;
     }
}

export class FullRecipeStepIngredient {
     id: number;
     name: string;
     ingredient: ValidIngredient;
     ingredientNotes: string;
     quantityType: QuantityType;
     quantityNotes: string;
     quantityValue: number;
     productOfRecipeStep: boolean;
     belongsToRecipeStep?: number;
     externalID: string;
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

     constructor(
         id: number = 0,
         name: string = '',
         ingredient: ValidIngredient = new ValidIngredient(),
         ingredientNotes: string = '',
         quantityType: QuantityType = 'grams',
         quantityNotes: string = '',
         quantityValue: number = 0,
         productOfRecipeStep: boolean = false,
         belongsToRecipeStep?: number,
         externalID: string = '',
         createdOn: number = 0,
         lastUpdatedOn?: number,
         archivedOn?: number,
     ) {
          this.id = id;
          this.name = name;
          this.ingredient = ingredient;
          this.ingredientNotes = ingredientNotes;
          this.quantityType = quantityType;
          this.quantityNotes = quantityNotes;
          this.quantityValue = quantityValue;
          this.productOfRecipeStep = productOfRecipeStep;
          this.belongsToRecipeStep = belongsToRecipeStep;
          this.externalID = externalID;
          this.createdOn = createdOn;
          this.lastUpdatedOn = lastUpdatedOn;
          this.archivedOn = archivedOn;
     }
}

export class RecipeStepIngredientList {
     recipeStepIngredients: RecipeStepIngredient[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          recipeStepIngredients: RecipeStepIngredient[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.recipeStepIngredients = recipeStepIngredients;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}