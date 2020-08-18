import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStepInstrument {
  id: number;
  instrumentID?: number;
  recipeStepID: number;
  notes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.instrumentID = 0;
    this.recipeStepID = 0;
    this.notes = "";
    this.createdOn = 0;
  }

static areEqual = function(
  rsi1: RecipeStepInstrument,
  rsi2: RecipeStepInstrument,
): boolean {
    return (
      rsi1.id === rsi2.id &&
      rsi1.instrumentID === rsi2.instrumentID &&
      rsi1.recipeStepID === rsi2.recipeStepID &&
      rsi1.notes === rsi2.notes &&
      rsi1.archivedOn === rsi2.archivedOn
    );
  }
}

export const fakeRecipeStepInstrumentFactory = Factory.Sync.makeFactory<RecipeStepInstrument> ({
  instrumentID: Factory.Sync.each(() =>  faker.random.number()),
  recipeStepID: Factory.Sync.each(() =>  faker.random.number()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
