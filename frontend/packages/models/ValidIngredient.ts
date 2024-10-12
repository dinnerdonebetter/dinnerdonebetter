// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRange } from './number_range';


export interface IValidIngredient {
   animalDerived: boolean;
 animalFlesh: boolean;
 archivedAt: string;
 containsAlcohol: boolean;
 containsDairy: boolean;
 containsEgg: boolean;
 containsFish: boolean;
 containsGluten: boolean;
 containsPeanut: boolean;
 containsSesame: boolean;
 containsShellfish: boolean;
 containsSoy: boolean;
 containsTreeNut: boolean;
 containsWheat: boolean;
 createdAt: string;
 description: string;
 iconPath: string;
 id: string;
 isAcid: boolean;
 isFat: boolean;
 isFruit: boolean;
 isGrain: boolean;
 isHeat: boolean;
 isLiquid: boolean;
 isProtein: boolean;
 isSalt: boolean;
 isStarch: boolean;
 lastUpdatedAt: string;
 name: string;
 pluralName: string;
 restrictToPreparations: boolean;
 shoppingSuggestions: string;
 slug: string;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 warning: string;

}

export class ValidIngredient implements IValidIngredient {
   animalDerived: boolean;
 animalFlesh: boolean;
 archivedAt: string;
 containsAlcohol: boolean;
 containsDairy: boolean;
 containsEgg: boolean;
 containsFish: boolean;
 containsGluten: boolean;
 containsPeanut: boolean;
 containsSesame: boolean;
 containsShellfish: boolean;
 containsSoy: boolean;
 containsTreeNut: boolean;
 containsWheat: boolean;
 createdAt: string;
 description: string;
 iconPath: string;
 id: string;
 isAcid: boolean;
 isFat: boolean;
 isFruit: boolean;
 isGrain: boolean;
 isHeat: boolean;
 isLiquid: boolean;
 isProtein: boolean;
 isSalt: boolean;
 isStarch: boolean;
 lastUpdatedAt: string;
 name: string;
 pluralName: string;
 restrictToPreparations: boolean;
 shoppingSuggestions: string;
 slug: string;
 storageInstructions: string;
 storageTemperatureInCelsius: NumberRange;
 warning: string;
constructor(input: Partial<ValidIngredient> = {}) {
	 this.animalDerived = input.animalDerived || false;
 this.animalFlesh = input.animalFlesh || false;
 this.archivedAt = input.archivedAt || '';
 this.containsAlcohol = input.containsAlcohol || false;
 this.containsDairy = input.containsDairy || false;
 this.containsEgg = input.containsEgg || false;
 this.containsFish = input.containsFish || false;
 this.containsGluten = input.containsGluten || false;
 this.containsPeanut = input.containsPeanut || false;
 this.containsSesame = input.containsSesame || false;
 this.containsShellfish = input.containsShellfish || false;
 this.containsSoy = input.containsSoy || false;
 this.containsTreeNut = input.containsTreeNut || false;
 this.containsWheat = input.containsWheat || false;
 this.createdAt = input.createdAt || '';
 this.description = input.description || '';
 this.iconPath = input.iconPath || '';
 this.id = input.id || '';
 this.isAcid = input.isAcid || false;
 this.isFat = input.isFat || false;
 this.isFruit = input.isFruit || false;
 this.isGrain = input.isGrain || false;
 this.isHeat = input.isHeat || false;
 this.isLiquid = input.isLiquid || false;
 this.isProtein = input.isProtein || false;
 this.isSalt = input.isSalt || false;
 this.isStarch = input.isStarch || false;
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.name = input.name || '';
 this.pluralName = input.pluralName || '';
 this.restrictToPreparations = input.restrictToPreparations || false;
 this.shoppingSuggestions = input.shoppingSuggestions || '';
 this.slug = input.slug || '';
 this.storageInstructions = input.storageInstructions || '';
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
 this.warning = input.warning || '';
}
}