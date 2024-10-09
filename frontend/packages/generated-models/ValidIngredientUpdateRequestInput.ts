// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRange } from './number_range';

export interface IValidIngredientUpdateRequestInput {
  containsSesame?: boolean;
  containsShellfish?: boolean;
  containsWheat?: boolean;
  description?: string;
  isGrain?: boolean;
  animalDerived?: boolean;
  containsPeanut?: boolean;
  storageInstructions?: string;
  isSalt?: boolean;
  restrictToPreparations?: boolean;
  animalFlesh?: boolean;
  isLiquid?: boolean;
  slug?: string;
  isAcid?: boolean;
  isFat?: boolean;
  isHeat?: boolean;
  containsDairy?: boolean;
  isFruit?: boolean;
  containsAlcohol?: boolean;
  containsSoy?: boolean;
  containsTreeNut?: boolean;
  iconPath?: string;
  isProtein?: boolean;
  name?: string;
  warning?: string;
  containsEgg?: boolean;
  containsGluten?: boolean;
  pluralName?: string;
  shoppingSuggestions?: string;
  storageTemperatureInCelsius: NumberRange;
  containsFish?: boolean;
  isStarch?: boolean;
}

export class ValidIngredientUpdateRequestInput implements IValidIngredientUpdateRequestInput {
  containsSesame?: boolean;
  containsShellfish?: boolean;
  containsWheat?: boolean;
  description?: string;
  isGrain?: boolean;
  animalDerived?: boolean;
  containsPeanut?: boolean;
  storageInstructions?: string;
  isSalt?: boolean;
  restrictToPreparations?: boolean;
  animalFlesh?: boolean;
  isLiquid?: boolean;
  slug?: string;
  isAcid?: boolean;
  isFat?: boolean;
  isHeat?: boolean;
  containsDairy?: boolean;
  isFruit?: boolean;
  containsAlcohol?: boolean;
  containsSoy?: boolean;
  containsTreeNut?: boolean;
  iconPath?: string;
  isProtein?: boolean;
  name?: string;
  warning?: string;
  containsEgg?: boolean;
  containsGluten?: boolean;
  pluralName?: string;
  shoppingSuggestions?: string;
  storageTemperatureInCelsius: NumberRange;
  containsFish?: boolean;
  isStarch?: boolean;
  constructor(input: Partial<ValidIngredientUpdateRequestInput> = {}) {
    this.containsSesame = input.containsSesame;
    this.containsShellfish = input.containsShellfish;
    this.containsWheat = input.containsWheat;
    this.description = input.description;
    this.isGrain = input.isGrain;
    this.animalDerived = input.animalDerived;
    this.containsPeanut = input.containsPeanut;
    this.storageInstructions = input.storageInstructions;
    this.isSalt = input.isSalt;
    this.restrictToPreparations = input.restrictToPreparations;
    this.animalFlesh = input.animalFlesh;
    this.isLiquid = input.isLiquid;
    this.slug = input.slug;
    this.isAcid = input.isAcid;
    this.isFat = input.isFat;
    this.isHeat = input.isHeat;
    this.containsDairy = input.containsDairy;
    this.isFruit = input.isFruit;
    this.containsAlcohol = input.containsAlcohol;
    this.containsSoy = input.containsSoy;
    this.containsTreeNut = input.containsTreeNut;
    this.iconPath = input.iconPath;
    this.isProtein = input.isProtein;
    this.name = input.name;
    this.warning = input.warning;
    this.containsEgg = input.containsEgg;
    this.containsGluten = input.containsGluten;
    this.pluralName = input.pluralName;
    this.shoppingSuggestions = input.shoppingSuggestions;
    this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
    this.containsFish = input.containsFish;
    this.isStarch = input.isStarch;
  }
}
