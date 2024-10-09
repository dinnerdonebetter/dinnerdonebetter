// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientUpdateRequestInput {
  containsShellfish?: boolean;
  description?: string;
  isAcid?: boolean;
  containsEgg?: boolean;
  containsTreeNut?: boolean;
  storageInstructions?: string;
  warning?: string;
  isStarch?: boolean;
  pluralName?: string;
  animalDerived?: boolean;
  containsAlcohol?: boolean;
  containsDairy?: boolean;
  isGrain?: boolean;
  isProtein?: boolean;
  containsSoy?: boolean;
  containsWheat?: boolean;
  isFruit?: boolean;
  restrictToPreparations?: boolean;
  storageTemperatureInCelsius: NumberRange;
  containsPeanut?: boolean;
  isHeat?: boolean;
  isSalt?: boolean;
  name?: string;
  slug?: string;
  shoppingSuggestions?: string;
  animalFlesh?: boolean;
  containsGluten?: boolean;
  iconPath?: string;
  isFat?: boolean;
  isLiquid?: boolean;
  containsFish?: boolean;
  containsSesame?: boolean;
}

export class ValidIngredientUpdateRequestInput implements IValidIngredientUpdateRequestInput {
  containsShellfish?: boolean;
  description?: string;
  isAcid?: boolean;
  containsEgg?: boolean;
  containsTreeNut?: boolean;
  storageInstructions?: string;
  warning?: string;
  isStarch?: boolean;
  pluralName?: string;
  animalDerived?: boolean;
  containsAlcohol?: boolean;
  containsDairy?: boolean;
  isGrain?: boolean;
  isProtein?: boolean;
  containsSoy?: boolean;
  containsWheat?: boolean;
  isFruit?: boolean;
  restrictToPreparations?: boolean;
  storageTemperatureInCelsius: NumberRange;
  containsPeanut?: boolean;
  isHeat?: boolean;
  isSalt?: boolean;
  name?: string;
  slug?: string;
  shoppingSuggestions?: string;
  animalFlesh?: boolean;
  containsGluten?: boolean;
  iconPath?: string;
  isFat?: boolean;
  isLiquid?: boolean;
  containsFish?: boolean;
  containsSesame?: boolean;
  constructor(input: Partial<ValidIngredientUpdateRequestInput> = {}) {
    this.containsShellfish = input.containsShellfish;
    this.description = input.description;
    this.isAcid = input.isAcid;
    this.containsEgg = input.containsEgg;
    this.containsTreeNut = input.containsTreeNut;
    this.storageInstructions = input.storageInstructions;
    this.warning = input.warning;
    this.isStarch = input.isStarch;
    this.pluralName = input.pluralName;
    this.animalDerived = input.animalDerived;
    this.containsAlcohol = input.containsAlcohol;
    this.containsDairy = input.containsDairy;
    this.isGrain = input.isGrain;
    this.isProtein = input.isProtein;
    this.containsSoy = input.containsSoy;
    this.containsWheat = input.containsWheat;
    this.isFruit = input.isFruit;
    this.restrictToPreparations = input.restrictToPreparations;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.containsPeanut = input.containsPeanut;
    this.isHeat = input.isHeat;
    this.isSalt = input.isSalt;
    this.name = input.name;
    this.slug = input.slug;
    this.shoppingSuggestions = input.shoppingSuggestions;
    this.animalFlesh = input.animalFlesh;
    this.containsGluten = input.containsGluten;
    this.iconPath = input.iconPath;
    this.isFat = input.isFat;
    this.isLiquid = input.isLiquid;
    this.containsFish = input.containsFish;
    this.containsSesame = input.containsSesame;
  }
}
