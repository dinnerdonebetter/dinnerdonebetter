import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class ValidIngredientPreparation {
  id: number;
  notes: string;
  validPreparationID: string;
  validIngredientID: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.notes = "";
    this.validPreparationID = "";
    this.validIngredientID = "";
    this.createdOn = 0;
  }

static areEqual = function(
  vip1: ValidIngredientPreparation,
  vip2: ValidIngredientPreparation,
): boolean {
    return (
      vip1.id === vip2.id &&
      vip1.notes === vip2.notes &&
      vip1.validPreparationID === vip2.validPreparationID &&
      vip1.validIngredientID === vip2.validIngredientID &&
      vip1.archivedOn === vip2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidIngredientPreparation> ({
  notes: Factory.Sync.each(() =>  faker.random.word()),
  validPreparationID: Factory.Sync.each(() =>  faker.random.word()),
  validIngredientID: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
