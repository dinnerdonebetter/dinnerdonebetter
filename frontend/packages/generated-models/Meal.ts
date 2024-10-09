// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealComponent } from './MealComponent';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IMeal {
   createdAt: string;
 estimatedPortions: NumberRangeWithOptionalMax;
 id: string;
 name: string;
 archivedAt?: string;
 components: MealComponent;
 createdByUser: string;
 description: string;
 eligibleForMealPlans: boolean;
 lastUpdatedAt?: string;

}

export class Meal implements IMeal {
   createdAt: string;
 estimatedPortions: NumberRangeWithOptionalMax;
 id: string;
 name: string;
 archivedAt?: string;
 components: MealComponent;
 createdByUser: string;
 description: string;
 eligibleForMealPlans: boolean;
 lastUpdatedAt?: string;
constructor(input: Partial<Meal> = {}) {
	 this.createdAt = input.createdAt = '';
 this.estimatedPortions = input.estimatedPortions = { min: 0 };
 this.id = input.id = '';
 this.name = input.name = '';
 this.archivedAt = input.archivedAt;
 this.components = input.components = new MealComponent();
 this.createdByUser = input.createdByUser = '';
 this.description = input.description = '';
 this.eligibleForMealPlans = input.eligibleForMealPlans = false;
 this.lastUpdatedAt = input.lastUpdatedAt;
}
}