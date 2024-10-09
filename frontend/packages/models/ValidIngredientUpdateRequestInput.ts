// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientUpdateRequestInput {
  containsWheat: boolean;
  isFruit: boolean;
  slug: string;
  storageTemperatureInCelsius: NumberRange;
  containsSoy: boolean;
  isGrain: boolean;
  shoppingSuggestions: string;
  containsFish: boolean;
  isFat: boolean;
  isStarch: boolean;
  storageInstructions: string;
  warning: string;
  animalDerived: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  isProtein: boolean;
  containsAlcohol: boolean;
  iconPath: string;
  isHeat: boolean;
  restrictToPreparations: boolean;
  containsShellfish: boolean;
  description: string;
  isAcid: boolean;
  isLiquid: boolean;
  animalFlesh: boolean;
  containsDairy: boolean;
  containsTreeNut: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  isSalt: boolean;
  name: string;
  pluralName: string;
}

export class ValidIngredientUpdateRequestInput implements IValidIngredientUpdateRequestInput {
  containsWheat: boolean;
  isFruit: boolean;
  slug: string;
  storageTemperatureInCelsius: NumberRange;
  containsSoy: boolean;
  isGrain: boolean;
  shoppingSuggestions: string;
  containsFish: boolean;
  isFat: boolean;
  isStarch: boolean;
  storageInstructions: string;
  warning: string;
  animalDerived: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  isProtein: boolean;
  containsAlcohol: boolean;
  iconPath: string;
  isHeat: boolean;
  restrictToPreparations: boolean;
  containsShellfish: boolean;
  description: string;
  isAcid: boolean;
  isLiquid: boolean;
  animalFlesh: boolean;
  containsDairy: boolean;
  containsTreeNut: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  isSalt: boolean;
  name: string;
  pluralName: string;
  constructor(input: Partial<ValidIngredientUpdateRequestInput> = {}) {
    this.containsWheat = input.containsWheat || false;
    this.isFruit = input.isFruit || false;
    this.slug = input.slug || '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.containsSoy = input.containsSoy || false;
    this.isGrain = input.isGrain || false;
    this.shoppingSuggestions = input.shoppingSuggestions || '';
    this.containsFish = input.containsFish || false;
    this.isFat = input.isFat || false;
    this.isStarch = input.isStarch || false;
    this.storageInstructions = input.storageInstructions || '';
    this.warning = input.warning || '';
    this.animalDerived = input.animalDerived || false;
    this.containsEgg = input.containsEgg || false;
    this.containsPeanut = input.containsPeanut || false;
    this.isProtein = input.isProtein || false;
    this.containsAlcohol = input.containsAlcohol || false;
    this.iconPath = input.iconPath || '';
    this.isHeat = input.isHeat || false;
    this.restrictToPreparations = input.restrictToPreparations || false;
    this.containsShellfish = input.containsShellfish || false;
    this.description = input.description || '';
    this.isAcid = input.isAcid || false;
    this.isLiquid = input.isLiquid || false;
    this.animalFlesh = input.animalFlesh || false;
    this.containsDairy = input.containsDairy || false;
    this.containsTreeNut = input.containsTreeNut || false;
    this.containsGluten = input.containsGluten || false;
    this.containsSesame = input.containsSesame || false;
    this.isSalt = input.isSalt || false;
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
  }
}
