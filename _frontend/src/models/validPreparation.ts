import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class ValidPreparation {
  id: number;
  name: string;
  description: string;
  icon: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = "";
    this.description = "";
    this.icon = "";
    this.createdOn = 0;
  }

static areEqual = function(
  vp1: ValidPreparation,
  vp2: ValidPreparation,
): boolean {
    return (
      vp1.id === vp2.id &&
      vp1.name === vp2.name &&
      vp1.description === vp2.description &&
      vp1.icon === vp2.icon &&
      vp1.archivedOn === vp2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidPreparation> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  description: Factory.Sync.each(() =>  faker.random.word()),
  icon: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
