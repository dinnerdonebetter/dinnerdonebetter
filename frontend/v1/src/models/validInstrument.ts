import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class ValidInstrument {
  id: number;
  name: string;
  variant: string;
  description: string;
  icon: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = "";
    this.variant = "";
    this.description = "";
    this.icon = "";
    this.createdOn = 0;
  }

static areEqual = function(
  vi1: ValidInstrument,
  vi2: ValidInstrument,
): boolean {
    return (
      vi1.id === vi2.id &&
      vi1.name === vi2.name &&
      vi1.variant === vi2.variant &&
      vi1.description === vi2.description &&
      vi1.icon === vi2.icon &&
      vi1.archivedOn === vi2.archivedOn
    );
  }
}

export const fakeValidInstrumentFactory = Factory.Sync.makeFactory<ValidInstrument> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  variant: Factory.Sync.each(() =>  faker.random.word()),
  description: Factory.Sync.each(() =>  faker.random.word()),
  icon: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
