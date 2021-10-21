import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class MealPlan {
  id: number;
  state: string;
  startsAt: number;
  endsAt: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.state = "";
    this.startsAt = 0;
    this.endsAt = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  mp1: MealPlan,
  mp2: MealPlan,
): boolean {
    return (
      mp1.id === mp2.id &&
      mp1.state === mp2.state &&
      mp1.startsAt === mp2.startsAt &&
      mp1.endsAt === mp2.endsAt &&
      mp1.archivedOn === mp2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<MealPlan> ({
  state: Factory.Sync.each(() =>  faker.random.word()),
  startsAt: Factory.Sync.each(() =>  faker.random.number()),
  endsAt: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
