// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientCreationRequestInput {
  containsDairy: boolean;
  isFruit: boolean;
  isLiquid: boolean;
  isStarch: boolean;
  containsSoy: boolean;
  containsWheat: boolean;
  storageInstructions: string;
  animalDerived: boolean;
  containsPeanut: boolean;
  iconPath: string;
  isAcid: boolean;
  slug: string;
  containsShellfish: boolean;
  description: string;
  isFat: boolean;
  isSalt: boolean;
  storageTemperatureInCelsius: NumberRange;
  warning: string;
  containsFish: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  isGrain: boolean;
  containsAlcohol: boolean;
  shoppingSuggestions: string;
  animalFlesh: boolean;
  containsEgg: boolean;
  isProtein: boolean;
  name: string;
  pluralName: string;
  restrictToPreparations: boolean;
  containsTreeNut: boolean;
  isHeat: boolean;
}

export class ValidIngredientCreationRequestInput implements IValidIngredientCreationRequestInput {
  containsDairy: boolean;
  isFruit: boolean;
  isLiquid: boolean;
  isStarch: boolean;
  containsSoy: boolean;
  containsWheat: boolean;
  storageInstructions: string;
  animalDerived: boolean;
  containsPeanut: boolean;
  iconPath: string;
  isAcid: boolean;
  slug: string;
  containsShellfish: boolean;
  description: string;
  isFat: boolean;
  isSalt: boolean;
  storageTemperatureInCelsius: NumberRange;
  warning: string;
  containsFish: boolean;
  containsGluten: boolean;
  containsSesame: boolean;
  isGrain: boolean;
  containsAlcohol: boolean;
  shoppingSuggestions: string;
  animalFlesh: boolean;
  containsEgg: boolean;
  isProtein: boolean;
  name: string;
  pluralName: string;
  restrictToPreparations: boolean;
  containsTreeNut: boolean;
  isHeat: boolean;
  constructor(input: Partial<ValidIngredientCreationRequestInput> = {}) {
    this.containsDairy = input.containsDairy || false;
    this.isFruit = input.isFruit || false;
    this.isLiquid = input.isLiquid || false;
    this.isStarch = input.isStarch || false;
    this.containsSoy = input.containsSoy || false;
    this.containsWheat = input.containsWheat || false;
    this.storageInstructions = input.storageInstructions || '';
    this.animalDerived = input.animalDerived || false;
    this.containsPeanut = input.containsPeanut || false;
    this.iconPath = input.iconPath || '';
    this.isAcid = input.isAcid || false;
    this.slug = input.slug || '';
    this.containsShellfish = input.containsShellfish || false;
    this.description = input.description || '';
    this.isFat = input.isFat || false;
    this.isSalt = input.isSalt || false;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
    this.warning = input.warning || '';
    this.containsFish = input.containsFish || false;
    this.containsGluten = input.containsGluten || false;
    this.containsSesame = input.containsSesame || false;
    this.isGrain = input.isGrain || false;
    this.containsAlcohol = input.containsAlcohol || false;
    this.shoppingSuggestions = input.shoppingSuggestions || '';
    this.animalFlesh = input.animalFlesh || false;
    this.containsEgg = input.containsEgg || false;
    this.isProtein = input.isProtein || false;
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.restrictToPreparations = input.restrictToPreparations || false;
    this.containsTreeNut = input.containsTreeNut || false;
    this.isHeat = input.isHeat || false;
  }
}
