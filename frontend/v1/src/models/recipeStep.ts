import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class RecipeStep {
  id: number;
  index: number;
  preparationID: number;
  prerequisiteStep: number;
  minEstimatedTimeInSeconds: number;
  maxEstimatedTimeInSeconds: number;
  temperatureInCelsius?: number;
  notes: string;
  recipeID: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.index = 0;
    this.preparationID = 0;
    this.prerequisiteStep = 0;
    this.minEstimatedTimeInSeconds = 0;
    this.maxEstimatedTimeInSeconds = 0;
    this.temperatureInCelsius = 0;
    this.notes = "";
    this.recipeID = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  rs1: RecipeStep,
  rs2: RecipeStep,
): boolean {
    return (
      rs1.id === rs2.id &&
      rs1.index === rs2.index &&
      rs1.preparationID === rs2.preparationID &&
      rs1.prerequisiteStep === rs2.prerequisiteStep &&
      rs1.minEstimatedTimeInSeconds === rs2.minEstimatedTimeInSeconds &&
      rs1.maxEstimatedTimeInSeconds === rs2.maxEstimatedTimeInSeconds &&
      rs1.temperatureInCelsius === rs2.temperatureInCelsius &&
      rs1.notes === rs2.notes &&
      rs1.recipeID === rs2.recipeID &&
      rs1.archivedOn === rs2.archivedOn
    );
  }
}

export const fakeRecipeStepFactory = Factory.Sync.makeFactory<RecipeStep> ({
  index: Factory.Sync.each(() =>  faker.random.number()),
  preparationID: Factory.Sync.each(() =>  faker.random.number()),
  prerequisiteStep: Factory.Sync.each(() =>  faker.random.number()),
  minEstimatedTimeInSeconds: Factory.Sync.each(() =>  faker.random.number()),
  maxEstimatedTimeInSeconds: Factory.Sync.each(() =>  faker.random.number()),
  temperatureInCelsius: Factory.Sync.each(() =>  faker.random.number()),
  notes: Factory.Sync.each(() =>  faker.random.word()),
  recipeID: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
