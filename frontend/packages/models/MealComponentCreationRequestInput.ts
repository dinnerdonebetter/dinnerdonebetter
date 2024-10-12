// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealComponentType } from './enums';


export interface IMealComponentCreationRequestInput {
   componentType: MealComponentType;
 recipeID: string;
 recipeScale: number;

}

export class MealComponentCreationRequestInput implements IMealComponentCreationRequestInput {
   componentType: MealComponentType;
 recipeID: string;
 recipeScale: number;
constructor(input: Partial<MealComponentCreationRequestInput> = {}) {
	 this.componentType = input.componentType || 'unspecified';
 this.recipeID = input.recipeID || '';
 this.recipeScale = input.recipeScale || 0;
}
}