// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRange } from './number_range';


export interface IValidIngredientCreationRequestInput {
   isFat: boolean;
 isProtein: boolean;
 name: string;
 containsPeanut: boolean;
 containsShellfish: boolean;
 containsTreeNut: boolean;
 slug: string;
 animalDerived: boolean;
 containsWheat: boolean;
 pluralName: string;
 storageTemperatureInCelsius: NumberRange;
 containsDairy: boolean;
 containsSesame: boolean;
 isHeat: boolean;
 warning: string;
 containsEgg: boolean;
 containsGluten: boolean;
 isSalt: boolean;
 restrictToPreparations: boolean;
 containsAlcohol: boolean;
 containsSoy: boolean;
 isFruit: boolean;
 shoppingSuggestions: string;
 storageInstructions: string;
 containsFish: boolean;
 description: string;
 isAcid: boolean;
 iconPath: string;
 isLiquid: boolean;
 animalFlesh: boolean;
 isGrain: boolean;
 isStarch: boolean;

}

export class ValidIngredientCreationRequestInput implements IValidIngredientCreationRequestInput {
   isFat: boolean;
 isProtein: boolean;
 name: string;
 containsPeanut: boolean;
 containsShellfish: boolean;
 containsTreeNut: boolean;
 slug: string;
 animalDerived: boolean;
 containsWheat: boolean;
 pluralName: string;
 storageTemperatureInCelsius: NumberRange;
 containsDairy: boolean;
 containsSesame: boolean;
 isHeat: boolean;
 warning: string;
 containsEgg: boolean;
 containsGluten: boolean;
 isSalt: boolean;
 restrictToPreparations: boolean;
 containsAlcohol: boolean;
 containsSoy: boolean;
 isFruit: boolean;
 shoppingSuggestions: string;
 storageInstructions: string;
 containsFish: boolean;
 description: string;
 isAcid: boolean;
 iconPath: string;
 isLiquid: boolean;
 animalFlesh: boolean;
 isGrain: boolean;
 isStarch: boolean;
constructor(input: Partial<ValidIngredientCreationRequestInput> = {}) {
	 this.isFat = input.isFat = false;
 this.isProtein = input.isProtein = false;
 this.name = input.name = '';
 this.containsPeanut = input.containsPeanut = false;
 this.containsShellfish = input.containsShellfish = false;
 this.containsTreeNut = input.containsTreeNut = false;
 this.slug = input.slug = '';
 this.animalDerived = input.animalDerived = false;
 this.containsWheat = input.containsWheat = false;
 this.pluralName = input.pluralName = '';
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.containsDairy = input.containsDairy = false;
 this.containsSesame = input.containsSesame = false;
 this.isHeat = input.isHeat = false;
 this.warning = input.warning = '';
 this.containsEgg = input.containsEgg = false;
 this.containsGluten = input.containsGluten = false;
 this.isSalt = input.isSalt = false;
 this.restrictToPreparations = input.restrictToPreparations = false;
 this.containsAlcohol = input.containsAlcohol = false;
 this.containsSoy = input.containsSoy = false;
 this.isFruit = input.isFruit = false;
 this.shoppingSuggestions = input.shoppingSuggestions = '';
 this.storageInstructions = input.storageInstructions = '';
 this.containsFish = input.containsFish = false;
 this.description = input.description = '';
 this.isAcid = input.isAcid = false;
 this.iconPath = input.iconPath = '';
 this.isLiquid = input.isLiquid = false;
 this.animalFlesh = input.animalFlesh = false;
 this.isGrain = input.isGrain = false;
 this.isStarch = input.isStarch = false;
}
}