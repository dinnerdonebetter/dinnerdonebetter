import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class IterationMedia {
  id: number;
  path: string;
  mimetype: string;
  recipeIterationID: number;
  recipeStepID?: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.path = "";
    this.mimetype = "";
    this.recipeIterationID = 0;
    this.recipeStepID = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  im1: IterationMedia,
  im2: IterationMedia,
): boolean {
    return (
      im1.id === im2.id &&
      im1.path === im2.path &&
      im1.mimetype === im2.mimetype &&
      im1.recipeIterationID === im2.recipeIterationID &&
      im1.recipeStepID === im2.recipeStepID &&
      im1.archivedOn === im2.archivedOn
    );
  }
}

export const fakeIterationMediaFactory = Factory.Sync.makeFactory<IterationMedia> ({
  path: Factory.Sync.each(() =>  faker.random.word()),
  mimetype: Factory.Sync.each(() =>  faker.random.word()),
  recipeIterationID: Factory.Sync.each(() =>  faker.random.number()),
  recipeStepID: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
