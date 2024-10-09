// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientCreationRequestInput {
  containsSoy: boolean;
  storageInstructions: string;
  isFruit: boolean;
  isProtein: boolean;
  isStarch: boolean;
  animalDerived: boolean;
  containsPeanut: boolean;
  isAcid: boolean;
  animalFlesh: boolean;
  containsEgg: boolean;
  containsTreeNut: boolean;
  pluralName: string;
  restrictToPreparations: boolean;
  iconPath: string;
  isGrain: boolean;
  slug: string;
  containsDairy: boolean;
  description: string;
  storageTemperatureInCelsius: NumberRange;
  containsWheat: boolean;
  isFat: boolean;
  isLiquid: boolean;
  isSalt: boolean;
  warning: string;
  containsAlcohol: boolean;
  containsFish: boolean;
  containsShellfish: boolean;
  name: string;
  shoppingSuggestions: string;
  containsGluten: boolean;
  containsSesame: boolean;
  isHeat: boolean;
}

export class ValidIngredientCreationRequestInput implements IValidIngredientCreationRequestInput {
  containsSoy: boolean;
  storageInstructions: string;
  isFruit: boolean;
  isProtein: boolean;
  isStarch: boolean;
  animalDerived: boolean;
  containsPeanut: boolean;
  isAcid: boolean;
  animalFlesh: boolean;
  containsEgg: boolean;
  containsTreeNut: boolean;
  pluralName: string;
  restrictToPreparations: boolean;
  iconPath: string;
  isGrain: boolean;
  slug: string;
  containsDairy: boolean;
  description: string;
  storageTemperatureInCelsius: NumberRange;
  containsWheat: boolean;
  isFat: boolean;
  isLiquid: boolean;
  isSalt: boolean;
  warning: string;
  containsAlcohol: boolean;
  containsFish: boolean;
  containsShellfish: boolean;
  name: string;
  shoppingSuggestions: string;
  containsGluten: boolean;
  containsSesame: boolean;
  isHeat: boolean;
  constructor(input: Partial<ValidIngredientCreationRequestInput> = {}) {
    this.containsSoy = input.containsSoy = false;
    this.storageInstructions = input.storageInstructions = '';
    this.isFruit = input.isFruit = false;
    this.isProtein = input.isProtein = false;
    this.isStarch = input.isStarch = false;
    this.animalDerived = input.animalDerived = false;
    this.containsPeanut = input.containsPeanut = false;
    this.isAcid = input.isAcid = false;
    this.animalFlesh = input.animalFlesh = false;
    this.containsEgg = input.containsEgg = false;
    this.containsTreeNut = input.containsTreeNut = false;
    this.pluralName = input.pluralName = '';
    this.restrictToPreparations = input.restrictToPreparations = false;
    this.iconPath = input.iconPath = '';
    this.isGrain = input.isGrain = false;
    this.slug = input.slug = '';
    this.containsDairy = input.containsDairy = false;
    this.description = input.description = '';
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.containsWheat = input.containsWheat = false;
    this.isFat = input.isFat = false;
    this.isLiquid = input.isLiquid = false;
    this.isSalt = input.isSalt = false;
    this.warning = input.warning = '';
    this.containsAlcohol = input.containsAlcohol = false;
    this.containsFish = input.containsFish = false;
    this.containsShellfish = input.containsShellfish = false;
    this.name = input.name = '';
    this.shoppingSuggestions = input.shoppingSuggestions = '';
    this.containsGluten = input.containsGluten = false;
    this.containsSesame = input.containsSesame = false;
    this.isHeat = input.isHeat = false;
  }
}
