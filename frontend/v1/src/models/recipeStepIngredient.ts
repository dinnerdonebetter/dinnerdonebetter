import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStepIngredient {
  id: number;
  ingredientID?: number;
  quantityType: string;
  quantityValue: number;
  quantityNotes: string;
  productOfRecipe: boolean;
  ingredientNotes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.ingredientID = 0;
    this.quantityType = "";
    this.quantityValue = 0;
    this.quantityNotes = "";
    this.productOfRecipe = false;
    this.ingredientNotes = "";
    this.createdOn = 0;
  }

static areEqual = function(
  rsi1: RecipeStepIngredient,
  rsi2: RecipeStepIngredient,
): boolean {
    return (
      rsi1.id === rsi2.id &&
      rsi1.ingredientID === rsi2.ingredientID &&
      rsi1.quantityType === rsi2.quantityType &&
      rsi1.quantityValue === rsi2.quantityValue &&
      rsi1.quantityNotes === rsi2.quantityNotes &&
      rsi1.productOfRecipe === rsi2.productOfRecipe &&
      rsi1.ingredientNotes === rsi2.ingredientNotes &&
      rsi1.archivedOn === rsi2.archivedOn
    );
  }
}

export const fakeRecipeStepIngredientFactory = Factory.Sync.makeFactory<RecipeStepIngredient> ({
  ingredientID: Factory.Sync.each(() =>  faker.random.number()),
  quantityType: Factory.Sync.each(() =>  faker.random.word()),
  quantityValue: Factory.Sync.each(() =>  faker.random.number()),
  quantityNotes: Factory.Sync.each(() =>  faker.random.word()),
  productOfRecipe: Factory.Sync.each(() =>  faker.random.boolean()),
  ingredientNotes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
