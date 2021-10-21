import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStepInstrument {
  id: number;
  instrumentID?: string;
  recipeStepID: string;
  notes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.instrumentID = "";
    this.recipeStepID = "";
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

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<RecipeStepInstrument> ({
  instrumentID: Factory.Sync.each(() =>  faker.random.word()),
  recipeStepID: Factory.Sync.each(() =>  faker.random.word()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
