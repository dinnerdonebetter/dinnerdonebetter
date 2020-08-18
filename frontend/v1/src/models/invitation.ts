import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class Invitation {
  id: number;
  code: string;
  consumed: boolean;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.code = "";
    this.consumed = false;
    this.createdOn = 0;
  }

static areEqual = function(
  i1: Invitation,
  i2: Invitation,
): boolean {
    return (
      i1.id === i2.id &&
      i1.code === i2.code &&
      i1.consumed === i2.consumed &&
      i1.archivedOn === i2.archivedOn
    );
  }
}

export const fakeInvitationFactory = Factory.Sync.makeFactory<Invitation> ({
  code: Factory.Sync.each(() =>  faker.random.word()),
  consumed: Factory.Sync.each(() =>  faker.random.boolean()),
  ...defaultFactories,
});
