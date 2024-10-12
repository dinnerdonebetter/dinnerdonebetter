// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealComponentCreationRequestInput } from './MealComponentCreationRequestInput';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IMealCreationRequestInput {
   components: MealComponentCreationRequestInput[];
 description: string;
 eligibleForMealPlans: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 name: string;

}

export class MealCreationRequestInput implements IMealCreationRequestInput {
   components: MealComponentCreationRequestInput[];
 description: string;
 eligibleForMealPlans: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 name: string;
constructor(input: Partial<MealCreationRequestInput> = {}) {
	 this.components = input.components || [];
 this.description = input.description || '';
 this.eligibleForMealPlans = input.eligibleForMealPlans || false;
 this.estimatedPortions = input.estimatedPortions || { min: 0 };
 this.name = input.name || '';
}
}