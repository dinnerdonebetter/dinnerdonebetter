// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientUpdateRequestInput {
  animalDerived?: boolean;
  containsPeanut?: boolean;
  isFat?: boolean;
  isHeat?: boolean;
  storageInstructions?: string;
  isStarch?: boolean;
  containsDairy?: boolean;
  containsFish?: boolean;
  containsSoy?: boolean;
  description?: string;
  containsTreeNut?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  containsShellfish?: boolean;
  iconPath?: string;
  isAcid?: boolean;
  isSalt?: boolean;
  warning?: string;
  containsEgg?: boolean;
  containsGluten?: boolean;
  isFruit?: boolean;
  restrictToPreparations?: boolean;
  shoppingSuggestions?: string;
  storageTemperatureInCelsius: NumberRange;
  animalFlesh?: boolean;
  containsAlcohol?: boolean;
  containsWheat?: boolean;
  isLiquid?: boolean;
  containsSesame?: boolean;
  isGrain?: boolean;
  isProtein?: boolean;
}

export class ValidIngredientUpdateRequestInput implements IValidIngredientUpdateRequestInput {
  animalDerived?: boolean;
  containsPeanut?: boolean;
  isFat?: boolean;
  isHeat?: boolean;
  storageInstructions?: string;
  isStarch?: boolean;
  containsDairy?: boolean;
  containsFish?: boolean;
  containsSoy?: boolean;
  description?: string;
  containsTreeNut?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  containsShellfish?: boolean;
  iconPath?: string;
  isAcid?: boolean;
  isSalt?: boolean;
  warning?: string;
  containsEgg?: boolean;
  containsGluten?: boolean;
  isFruit?: boolean;
  restrictToPreparations?: boolean;
  shoppingSuggestions?: string;
  storageTemperatureInCelsius: NumberRange;
  animalFlesh?: boolean;
  containsAlcohol?: boolean;
  containsWheat?: boolean;
  isLiquid?: boolean;
  containsSesame?: boolean;
  isGrain?: boolean;
  isProtein?: boolean;
  constructor(input: Partial<ValidIngredientUpdateRequestInput> = {}) {
    this.animalDerived = input.animalDerived;
    this.containsPeanut = input.containsPeanut;
    this.isFat = input.isFat;
    this.isHeat = input.isHeat;
    this.storageInstructions = input.storageInstructions;
    this.isStarch = input.isStarch;
    this.containsDairy = input.containsDairy;
    this.containsFish = input.containsFish;
    this.containsSoy = input.containsSoy;
    this.description = input.description;
    this.containsTreeNut = input.containsTreeNut;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.containsShellfish = input.containsShellfish;
    this.iconPath = input.iconPath;
    this.isAcid = input.isAcid;
    this.isSalt = input.isSalt;
    this.warning = input.warning;
    this.containsEgg = input.containsEgg;
    this.containsGluten = input.containsGluten;
    this.isFruit = input.isFruit;
    this.restrictToPreparations = input.restrictToPreparations;
    this.shoppingSuggestions = input.shoppingSuggestions;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.animalFlesh = input.animalFlesh;
    this.containsAlcohol = input.containsAlcohol;
    this.containsWheat = input.containsWheat;
    this.isLiquid = input.isLiquid;
    this.containsSesame = input.containsSesame;
    this.isGrain = input.isGrain;
    this.isProtein = input.isProtein;
  }
}
