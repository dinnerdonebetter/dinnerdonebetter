// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredient {
  containsFish: boolean;
  isFruit: boolean;
  isGrain: boolean;
  animalFlesh: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  containsShellfish: boolean;
  isAcid: boolean;
  isHeat: boolean;
  warning: string;
  archivedAt?: string;
  containsTreeNut: boolean;
  restrictToPreparations: boolean;
  containsDairy: boolean;
  containsWheat: boolean;
  createdAt: string;
  isLiquid: boolean;
  animalDerived: boolean;
  description: string;
  isProtein: boolean;
  isSalt: boolean;
  iconPath: string;
  name: string;
  pluralName: string;
  slug: string;
  containsAlcohol: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  containsSoy: boolean;
  isFat: boolean;
  isStarch: boolean;
  shoppingSuggestions: string;
  id: string;
  lastUpdatedAt?: string;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
}

export class ValidIngredient implements IValidIngredient {
  containsFish: boolean;
  isFruit: boolean;
  isGrain: boolean;
  animalFlesh: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  containsShellfish: boolean;
  isAcid: boolean;
  isHeat: boolean;
  warning: string;
  archivedAt?: string;
  containsTreeNut: boolean;
  restrictToPreparations: boolean;
  containsDairy: boolean;
  containsWheat: boolean;
  createdAt: string;
  isLiquid: boolean;
  animalDerived: boolean;
  description: string;
  isProtein: boolean;
  isSalt: boolean;
  iconPath: string;
  name: string;
  pluralName: string;
  slug: string;
  containsAlcohol: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  containsSoy: boolean;
  isFat: boolean;
  isStarch: boolean;
  shoppingSuggestions: string;
  id: string;
  lastUpdatedAt?: string;
  storageInstructions: string;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<ValidIngredient> = {}) {
    this.containsFish = input.containsFish = false;
    this.isFruit = input.isFruit = false;
    this.isGrain = input.isGrain = false;
    this.animalFlesh = input.animalFlesh = false;
    this.containsGluten = input.containsGluten = false;
    this.containsSesame = input.containsSesame = false;
    this.containsShellfish = input.containsShellfish = false;
    this.isAcid = input.isAcid = false;
    this.isHeat = input.isHeat = false;
    this.warning = input.warning = '';
    this.archivedAt = input.archivedAt;
    this.containsTreeNut = input.containsTreeNut = false;
    this.restrictToPreparations = input.restrictToPreparations = false;
    this.containsDairy = input.containsDairy = false;
    this.containsWheat = input.containsWheat = false;
    this.createdAt = input.createdAt = '';
    this.isLiquid = input.isLiquid = false;
    this.animalDerived = input.animalDerived = false;
    this.description = input.description = '';
    this.isProtein = input.isProtein = false;
    this.isSalt = input.isSalt = false;
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.containsAlcohol = input.containsAlcohol = false;
    this.containsEgg = input.containsEgg = false;
    this.containsPeanut = input.containsPeanut = false;
    this.containsSoy = input.containsSoy = false;
    this.isFat = input.isFat = false;
    this.isStarch = input.isStarch = false;
    this.shoppingSuggestions = input.shoppingSuggestions = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.storageInstructions = input.storageInstructions = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
  }
}
