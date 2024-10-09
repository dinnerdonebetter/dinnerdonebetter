// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRange } from './number_range';


export interface IValidIngredient {
   isStarch: boolean;
 lastUpdatedAt?: string;
 animalFlesh: boolean;
 containsFish: boolean;
 containsGluten: boolean;
 containsPeanut: boolean;
 isHeat: boolean;
 containsSesame: boolean;
 name: string;
 slug: string;
 id: string;
 isLiquid: boolean;
 restrictToPreparations: boolean;
 storageTemperatureInCelsius: NumberRange;
 isGrain: boolean;
 isProtein: boolean;
 shoppingSuggestions: string;
 containsEgg: boolean;
 storageInstructions: string;
 warning: string;
 isFat: boolean;
 animalDerived: boolean;
 archivedAt?: string;
 containsAlcohol: boolean;
 containsShellfish: boolean;
 containsSoy: boolean;
 iconPath: string;
 isAcid: boolean;
 isFruit: boolean;
 containsDairy: boolean;
 containsTreeNut: boolean;
 containsWheat: boolean;
 createdAt: string;
 description: string;
 isSalt: boolean;
 pluralName: string;

}

export class ValidIngredient implements IValidIngredient {
   isStarch: boolean;
 lastUpdatedAt?: string;
 animalFlesh: boolean;
 containsFish: boolean;
 containsGluten: boolean;
 containsPeanut: boolean;
 isHeat: boolean;
 containsSesame: boolean;
 name: string;
 slug: string;
 id: string;
 isLiquid: boolean;
 restrictToPreparations: boolean;
 storageTemperatureInCelsius: NumberRange;
 isGrain: boolean;
 isProtein: boolean;
 shoppingSuggestions: string;
 containsEgg: boolean;
 storageInstructions: string;
 warning: string;
 isFat: boolean;
 animalDerived: boolean;
 archivedAt?: string;
 containsAlcohol: boolean;
 containsShellfish: boolean;
 containsSoy: boolean;
 iconPath: string;
 isAcid: boolean;
 isFruit: boolean;
 containsDairy: boolean;
 containsTreeNut: boolean;
 containsWheat: boolean;
 createdAt: string;
 description: string;
 isSalt: boolean;
 pluralName: string;
constructor(input: Partial<ValidIngredient> = {}) {
	 this.isStarch = input.isStarch = false;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.animalFlesh = input.animalFlesh = false;
 this.containsFish = input.containsFish = false;
 this.containsGluten = input.containsGluten = false;
 this.containsPeanut = input.containsPeanut = false;
 this.isHeat = input.isHeat = false;
 this.containsSesame = input.containsSesame = false;
 this.name = input.name = '';
 this.slug = input.slug = '';
 this.id = input.id = '';
 this.isLiquid = input.isLiquid = false;
 this.restrictToPreparations = input.restrictToPreparations = false;
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.isGrain = input.isGrain = false;
 this.isProtein = input.isProtein = false;
 this.shoppingSuggestions = input.shoppingSuggestions = '';
 this.containsEgg = input.containsEgg = false;
 this.storageInstructions = input.storageInstructions = '';
 this.warning = input.warning = '';
 this.isFat = input.isFat = false;
 this.animalDerived = input.animalDerived = false;
 this.archivedAt = input.archivedAt;
 this.containsAlcohol = input.containsAlcohol = false;
 this.containsShellfish = input.containsShellfish = false;
 this.containsSoy = input.containsSoy = false;
 this.iconPath = input.iconPath = '';
 this.isAcid = input.isAcid = false;
 this.isFruit = input.isFruit = false;
 this.containsDairy = input.containsDairy = false;
 this.containsTreeNut = input.containsTreeNut = false;
 this.containsWheat = input.containsWheat = false;
 this.createdAt = input.createdAt = '';
 this.description = input.description = '';
 this.isSalt = input.isSalt = false;
 this.pluralName = input.pluralName = '';
}
}