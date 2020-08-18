import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RequiredPreparationInstrument {
  id: number;
  instrumentID: number;
  preparationID: number;
  notes: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.instrumentID = 0;
    this.preparationID = 0;
    this.notes = "";
    this.createdOn = 0;
  }

static areEqual = function(
  rpi1: RequiredPreparationInstrument,
  rpi2: RequiredPreparationInstrument,
): boolean {
    return (
      rpi1.id === rpi2.id &&
      rpi1.instrumentID === rpi2.instrumentID &&
      rpi1.preparationID === rpi2.preparationID &&
      rpi1.notes === rpi2.notes &&
      rpi1.archivedOn === rpi2.archivedOn
    );
  }
}

export const fakeRequiredPreparationInstrumentFactory = Factory.Sync.makeFactory<RequiredPreparationInstrument> ({
  instrumentID: Factory.Sync.each(() =>  faker.random.number()),
  preparationID: Factory.Sync.each(() =>  faker.random.number()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
