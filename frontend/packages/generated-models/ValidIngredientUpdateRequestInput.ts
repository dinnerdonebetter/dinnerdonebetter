// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRange } from './number_range';


export interface IValidIngredientUpdateRequestInput {
   containsPeanut?: boolean;
 description?: string;
 storageInstructions?: string;
 containsFish?: boolean;
 containsSesame?: boolean;
 containsTreeNut?: boolean;
 isLiquid?: boolean;
 isStarch?: boolean;
 restrictToPreparations?: boolean;
 containsAlcohol?: boolean;
 containsWheat?: boolean;
 iconPath?: string;
 isAcid?: boolean;
 isFat?: boolean;
 isProtein?: boolean;
 isSalt?: boolean;
 shoppingSuggestions?: string;
 containsDairy?: boolean;
 containsSoy?: boolean;
 isFruit?: boolean;
 isGrain?: boolean;
 warning?: string;
 animalFlesh?: boolean;
 containsGluten?: boolean;
 name?: string;
 animalDerived?: boolean;
 containsEgg?: boolean;
 containsShellfish?: boolean;
 isHeat?: boolean;
 pluralName?: string;
 slug?: string;
 storageTemperatureInCelsius: NumberRange;

}

export class ValidIngredientUpdateRequestInput implements IValidIngredientUpdateRequestInput {
   containsPeanut?: boolean;
 description?: string;
 storageInstructions?: string;
 containsFish?: boolean;
 containsSesame?: boolean;
 containsTreeNut?: boolean;
 isLiquid?: boolean;
 isStarch?: boolean;
 restrictToPreparations?: boolean;
 containsAlcohol?: boolean;
 containsWheat?: boolean;
 iconPath?: string;
 isAcid?: boolean;
 isFat?: boolean;
 isProtein?: boolean;
 isSalt?: boolean;
 shoppingSuggestions?: string;
 containsDairy?: boolean;
 containsSoy?: boolean;
 isFruit?: boolean;
 isGrain?: boolean;
 warning?: string;
 animalFlesh?: boolean;
 containsGluten?: boolean;
 name?: string;
 animalDerived?: boolean;
 containsEgg?: boolean;
 containsShellfish?: boolean;
 isHeat?: boolean;
 pluralName?: string;
 slug?: string;
 storageTemperatureInCelsius: NumberRange;
constructor(input: Partial<ValidIngredientUpdateRequestInput> = {}) {
	 this.containsPeanut = input.containsPeanut;
 this.description = input.description;
 this.storageInstructions = input.storageInstructions;
 this.containsFish = input.containsFish;
 this.containsSesame = input.containsSesame;
 this.containsTreeNut = input.containsTreeNut;
 this.isLiquid = input.isLiquid;
 this.isStarch = input.isStarch;
 this.restrictToPreparations = input.restrictToPreparations;
 this.containsAlcohol = input.containsAlcohol;
 this.containsWheat = input.containsWheat;
 this.iconPath = input.iconPath;
 this.isAcid = input.isAcid;
 this.isFat = input.isFat;
 this.isProtein = input.isProtein;
 this.isSalt = input.isSalt;
 this.shoppingSuggestions = input.shoppingSuggestions;
 this.containsDairy = input.containsDairy;
 this.containsSoy = input.containsSoy;
 this.isFruit = input.isFruit;
 this.isGrain = input.isGrain;
 this.warning = input.warning;
 this.animalFlesh = input.animalFlesh;
 this.containsGluten = input.containsGluten;
 this.name = input.name;
 this.animalDerived = input.animalDerived;
 this.containsEgg = input.containsEgg;
 this.containsShellfish = input.containsShellfish;
 this.isHeat = input.isHeat;
 this.pluralName = input.pluralName;
 this.slug = input.slug;
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
}
}