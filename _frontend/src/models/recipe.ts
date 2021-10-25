import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class Recipe {
  id: number;
  name: string;
  source: string;
  description: string;
  inspiredByRecipeID?: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = "";
    this.source = "";
    this.description = "";
    this.inspiredByRecipeID = "";
    this.createdOn = 0;
  }

static areEqual = function(
  r1: Recipe,
  r2: Recipe,
): boolean {
    return (
      r1.id === r2.id &&
      r1.name === r2.name &&
      r1.source === r2.source &&
      r1.description === r2.description &&
      r1.inspiredByRecipeID === r2.inspiredByRecipeID &&
      r1.archivedOn === r2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<Recipe> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  source: Factory.Sync.each(() =>  faker.random.word()),
  description: Factory.Sync.each(() =>  faker.random.word()),
  inspiredByRecipeID: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
