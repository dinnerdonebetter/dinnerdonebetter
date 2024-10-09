// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientCreationRequestInput {
  containsShellfish: boolean;
  description: string;
  slug: string;
  storageInstructions: string;
  animalDerived: boolean;
  animalFlesh: boolean;
  containsFish: boolean;
  isFruit: boolean;
  isHeat: boolean;
  pluralName: string;
  containsAlcohol: boolean;
  isProtein: boolean;
  isStarch: boolean;
  name: string;
  shoppingSuggestions: string;
  containsGluten: boolean;
  containsSesame: boolean;
  containsSoy: boolean;
  isGrain: boolean;
  isSalt: boolean;
  containsDairy: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  iconPath: string;
  warning: string;
  containsTreeNut: boolean;
  isAcid: boolean;
  isFat: boolean;
  isLiquid: boolean;
  restrictToPreparations: boolean;
  containsWheat: boolean;
  storageTemperatureInCelsius: NumberRange;
}

export class ValidIngredientCreationRequestInput implements IValidIngredientCreationRequestInput {
  containsShellfish: boolean;
  description: string;
  slug: string;
  storageInstructions: string;
  animalDerived: boolean;
  animalFlesh: boolean;
  containsFish: boolean;
  isFruit: boolean;
  isHeat: boolean;
  pluralName: string;
  containsAlcohol: boolean;
  isProtein: boolean;
  isStarch: boolean;
  name: string;
  shoppingSuggestions: string;
  containsGluten: boolean;
  containsSesame: boolean;
  containsSoy: boolean;
  isGrain: boolean;
  isSalt: boolean;
  containsDairy: boolean;
  containsEgg: boolean;
  containsPeanut: boolean;
  iconPath: string;
  warning: string;
  containsTreeNut: boolean;
  isAcid: boolean;
  isFat: boolean;
  isLiquid: boolean;
  restrictToPreparations: boolean;
  containsWheat: boolean;
  storageTemperatureInCelsius: NumberRange;
  constructor(input: Partial<ValidIngredientCreationRequestInput> = {}) {
    this.containsShellfish = input.containsShellfish = false;
    this.description = input.description = '';
    this.slug = input.slug = '';
    this.storageInstructions = input.storageInstructions = '';
    this.animalDerived = input.animalDerived = false;
    this.animalFlesh = input.animalFlesh = false;
    this.containsFish = input.containsFish = false;
    this.isFruit = input.isFruit = false;
    this.isHeat = input.isHeat = false;
    this.pluralName = input.pluralName = '';
    this.containsAlcohol = input.containsAlcohol = false;
    this.isProtein = input.isProtein = false;
    this.isStarch = input.isStarch = false;
    this.name = input.name = '';
    this.shoppingSuggestions = input.shoppingSuggestions = '';
    this.containsGluten = input.containsGluten = false;
    this.containsSesame = input.containsSesame = false;
    this.containsSoy = input.containsSoy = false;
    this.isGrain = input.isGrain = false;
    this.isSalt = input.isSalt = false;
    this.containsDairy = input.containsDairy = false;
    this.containsEgg = input.containsEgg = false;
    this.containsPeanut = input.containsPeanut = false;
    this.iconPath = input.iconPath = '';
    this.warning = input.warning = '';
    this.containsTreeNut = input.containsTreeNut = false;
    this.isAcid = input.isAcid = false;
    this.isFat = input.isFat = false;
    this.isLiquid = input.isLiquid = false;
    this.restrictToPreparations = input.restrictToPreparations = false;
    this.containsWheat = input.containsWheat = false;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
  }
}
