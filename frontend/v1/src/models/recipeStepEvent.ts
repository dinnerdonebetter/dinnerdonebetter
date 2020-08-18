import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStepEvent {
  id: number;
  eventType: string;
  done: boolean;
  recipeIterationID: number;
  recipeStepID: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.eventType = "";
    this.done = false;
    this.recipeIterationID = 0;
    this.recipeStepID = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  rse1: RecipeStepEvent,
  rse2: RecipeStepEvent,
): boolean {
    return (
      rse1.id === rse2.id &&
      rse1.eventType === rse2.eventType &&
      rse1.done === rse2.done &&
      rse1.recipeIterationID === rse2.recipeIterationID &&
      rse1.recipeStepID === rse2.recipeStepID &&
      rse1.archivedOn === rse2.archivedOn
    );
  }
}

export const fakeRecipeStepEventFactory = Factory.Sync.makeFactory<RecipeStepEvent> ({
  eventType: Factory.Sync.each(() =>  faker.random.word()),
  done: Factory.Sync.each(() =>  faker.random.boolean()),
  recipeIterationID: Factory.Sync.each(() =>  faker.random.number()),
  recipeStepID: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
