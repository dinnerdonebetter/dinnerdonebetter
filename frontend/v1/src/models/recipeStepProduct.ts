import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStepProduct {
  id: number;
  name: string;
  recipeStepID: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = "";
    this.recipeStepID = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  rsp1: RecipeStepProduct,
  rsp2: RecipeStepProduct,
): boolean {
    return (
      rsp1.id === rsp2.id &&
      rsp1.name === rsp2.name &&
      rsp1.recipeStepID === rsp2.recipeStepID &&
      rsp1.archivedOn === rsp2.archivedOn
    );
  }
}

export const fakeRecipeStepProductFactory = Factory.Sync.makeFactory<RecipeStepProduct> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  recipeStepID: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
