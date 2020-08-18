import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class ValidIngredient {
  id: number;
  name: string;
  variant: string;
  description: string;
  warning: string;
  containsEgg: boolean;
  containsDairy: boolean;
  containsPeanut: boolean;
  containsTreeNut: boolean;
  containsSoy: boolean;
  containsWheat: boolean;
  containsShellfish: boolean;
  containsSesame: boolean;
  containsFish: boolean;
  containsGluten: boolean;
  animalFlesh: boolean;
  animalDerived: boolean;
  measurableByVolume: boolean;
  icon: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = "";
    this.variant = "";
    this.description = "";
    this.warning = "";
    this.containsEgg = false;
    this.containsDairy = false;
    this.containsPeanut = false;
    this.containsTreeNut = false;
    this.containsSoy = false;
    this.containsWheat = false;
    this.containsShellfish = false;
    this.containsSesame = false;
    this.containsFish = false;
    this.containsGluten = false;
    this.animalFlesh = false;
    this.animalDerived = false;
    this.measurableByVolume = false;
    this.icon = "";
    this.createdOn = 0;
  }

static areEqual = function(
  vi1: ValidIngredient,
  vi2: ValidIngredient,
): boolean {
    return (
      vi1.id === vi2.id &&
      vi1.name === vi2.name &&
      vi1.variant === vi2.variant &&
      vi1.description === vi2.description &&
      vi1.warning === vi2.warning &&
      vi1.containsEgg === vi2.containsEgg &&
      vi1.containsDairy === vi2.containsDairy &&
      vi1.containsPeanut === vi2.containsPeanut &&
      vi1.containsTreeNut === vi2.containsTreeNut &&
      vi1.containsSoy === vi2.containsSoy &&
      vi1.containsWheat === vi2.containsWheat &&
      vi1.containsShellfish === vi2.containsShellfish &&
      vi1.containsSesame === vi2.containsSesame &&
      vi1.containsFish === vi2.containsFish &&
      vi1.containsGluten === vi2.containsGluten &&
      vi1.animalFlesh === vi2.animalFlesh &&
      vi1.animalDerived === vi2.animalDerived &&
      vi1.measurableByVolume === vi2.measurableByVolume &&
      vi1.icon === vi2.icon &&
      vi1.archivedOn === vi2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidIngredient> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  variant: Factory.Sync.each(() =>  faker.random.word()),
  description: Factory.Sync.each(() =>  faker.random.word()),
  warning: Factory.Sync.each(() =>  faker.random.word()),
  containsEgg: Factory.Sync.each(() =>  faker.random.boolean()),
  containsDairy: Factory.Sync.each(() =>  faker.random.boolean()),
  containsPeanut: Factory.Sync.each(() =>  faker.random.boolean()),
  containsTreeNut: Factory.Sync.each(() =>  faker.random.boolean()),
  containsSoy: Factory.Sync.each(() =>  faker.random.boolean()),
  containsWheat: Factory.Sync.each(() =>  faker.random.boolean()),
  containsShellfish: Factory.Sync.each(() =>  faker.random.boolean()),
  containsSesame: Factory.Sync.each(() =>  faker.random.boolean()),
  containsFish: Factory.Sync.each(() =>  faker.random.boolean()),
  containsGluten: Factory.Sync.each(() =>  faker.random.boolean()),
  animalFlesh: Factory.Sync.each(() =>  faker.random.boolean()),
  animalDerived: Factory.Sync.each(() =>  faker.random.boolean()),
  measurableByVolume: Factory.Sync.each(() =>  faker.random.boolean()),
  icon: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
