import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class MealPlanOptionVote {
  id: number;
  points: number;
  abstain: boolean;
  notes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.points = 0;
    this.abstain = false;
    this.notes = "";
    this.createdOn = 0;
  }

static areEqual = function(
  mpov1: MealPlanOptionVote,
  mpov2: MealPlanOptionVote,
): boolean {
    return (
      mpov1.id === mpov2.id &&
      mpov1.points === mpov2.points &&
      mpov1.abstain === mpov2.abstain &&
      mpov1.notes === mpov2.notes &&
      mpov1.archivedOn === mpov2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<MealPlanOptionVote> ({
  points: Factory.Sync.each(() =>  faker.random.number()),
  abstain: Factory.Sync.each(() =>  faker.random.boolean()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
