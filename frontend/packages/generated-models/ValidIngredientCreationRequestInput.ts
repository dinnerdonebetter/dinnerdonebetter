// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientCreationRequestInput {
  containsTreeNut: boolean;
  isFruit: boolean;
  isGrain: boolean;
  pluralName: string;
  containsEgg: boolean;
  containsSoy: boolean;
  isHeat: boolean;
  name: string;
  storageInstructions: string;
  description: string;
  isAcid: boolean;
  restrictToPreparations: boolean;
  containsWheat: boolean;
  warning: string;
  containsShellfish: boolean;
  slug: string;
  containsAlcohol: boolean;
  containsPeanut: boolean;
  isLiquid: boolean;
  isStarch: boolean;
  containsDairy: boolean;
  containsFish: boolean;
  containsSesame: boolean;
  isFat: boolean;
  isProtein: boolean;
  isSalt: boolean;
  storageTemperatureInCelsius: NumberRange;
  animalDerived: boolean;
  animalFlesh: boolean;
  containsGluten: boolean;
  iconPath: string;
  shoppingSuggestions: string;
}

export class ValidIngredientCreationRequestInput implements IValidIngredientCreationRequestInput {
  containsTreeNut: boolean;
  isFruit: boolean;
  isGrain: boolean;
  pluralName: string;
  containsEgg: boolean;
  containsSoy: boolean;
  isHeat: boolean;
  name: string;
  storageInstructions: string;
  description: string;
  isAcid: boolean;
  restrictToPreparations: boolean;
  containsWheat: boolean;
  warning: string;
  containsShellfish: boolean;
  slug: string;
  containsAlcohol: boolean;
  containsPeanut: boolean;
  isLiquid: boolean;
  isStarch: boolean;
  containsDairy: boolean;
  containsFish: boolean;
  containsSesame: boolean;
  isFat: boolean;
  isProtein: boolean;
  isSalt: boolean;
  storageTemperatureInCelsius: NumberRange;
  animalDerived: boolean;
  animalFlesh: boolean;
  containsGluten: boolean;
  iconPath: string;
  shoppingSuggestions: string;
  constructor(input: Partial<ValidIngredientCreationRequestInput> = {}) {
    this.containsTreeNut = input.containsTreeNut = false;
    this.isFruit = input.isFruit = false;
    this.isGrain = input.isGrain = false;
    this.pluralName = input.pluralName = '';
    this.containsEgg = input.containsEgg = false;
    this.containsSoy = input.containsSoy = false;
    this.isHeat = input.isHeat = false;
    this.name = input.name = '';
    this.storageInstructions = input.storageInstructions = '';
    this.description = input.description = '';
    this.isAcid = input.isAcid = false;
    this.restrictToPreparations = input.restrictToPreparations = false;
    this.containsWheat = input.containsWheat = false;
    this.warning = input.warning = '';
    this.containsShellfish = input.containsShellfish = false;
    this.slug = input.slug = '';
    this.containsAlcohol = input.containsAlcohol = false;
    this.containsPeanut = input.containsPeanut = false;
    this.isLiquid = input.isLiquid = false;
    this.isStarch = input.isStarch = false;
    this.containsDairy = input.containsDairy = false;
    this.containsFish = input.containsFish = false;
    this.containsSesame = input.containsSesame = false;
    this.isFat = input.isFat = false;
    this.isProtein = input.isProtein = false;
    this.isSalt = input.isSalt = false;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.animalDerived = input.animalDerived = false;
    this.animalFlesh = input.animalFlesh = false;
    this.containsGluten = input.containsGluten = false;
    this.iconPath = input.iconPath = '';
    this.shoppingSuggestions = input.shoppingSuggestions = '';
  }
}
