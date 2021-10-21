import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class MealPlanOption {
  id: number;
  mealPlanID: string;
  dayOfWeek: number;
  recipeID: string;
  notes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.mealPlanID = "";
    this.dayOfWeek = 0;
    this.recipeID = "";
    this.notes = "";
    this.createdOn = 0;
  }

static areEqual = function(
  mpo1: MealPlanOption,
  mpo2: MealPlanOption,
): boolean {
    return (
      mpo1.id === mpo2.id &&
      mpo1.mealPlanID === mpo2.mealPlanID &&
      mpo1.dayOfWeek === mpo2.dayOfWeek &&
      mpo1.recipeID === mpo2.recipeID &&
      mpo1.notes === mpo2.notes &&
      mpo1.archivedOn === mpo2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<MealPlanOption> ({
  mealPlanID: Factory.Sync.each(() =>  faker.random.word()),
  dayOfWeek: Factory.Sync.each(() =>  faker.random.number()),
  recipeID: Factory.Sync.each(() =>  faker.random.word()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
