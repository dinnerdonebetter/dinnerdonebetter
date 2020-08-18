import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeIteration {
  id: number;
  recipeID: number;
  endDifficultyRating: number;
  endComplexityRating: number;
  endTasteRating: number;
  endOverallRating: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.recipeID = 0;
    this.endDifficultyRating = 0;
    this.endComplexityRating = 0;
    this.endTasteRating = 0;
    this.endOverallRating = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  ri1: RecipeIteration,
  ri2: RecipeIteration,
): boolean {
    return (
      ri1.id === ri2.id &&
      ri1.recipeID === ri2.recipeID &&
      ri1.endDifficultyRating === ri2.endDifficultyRating &&
      ri1.endComplexityRating === ri2.endComplexityRating &&
      ri1.endTasteRating === ri2.endTasteRating &&
      ri1.endOverallRating === ri2.endOverallRating &&
      ri1.archivedOn === ri2.archivedOn
    );
  }
}

export const fakeRecipeIterationFactory = Factory.Sync.makeFactory<RecipeIteration> ({
  recipeID: Factory.Sync.each(() =>  faker.random.number()),
  endDifficultyRating: Factory.Sync.each(() =>  faker.random.number()),
  endComplexityRating: Factory.Sync.each(() =>  faker.random.number()),
  endTasteRating: Factory.Sync.each(() =>  faker.random.number()),
  endOverallRating: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
