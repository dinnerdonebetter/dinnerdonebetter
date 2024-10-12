// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealComponent } from './MealComponent';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IMeal {
   archivedAt: string;
 components: MealComponent[];
 createdAt: string;
 createdByUser: string;
 description: string;
 eligibleForMealPlans: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 id: string;
 lastUpdatedAt: string;
 name: string;

}

export class Meal implements IMeal {
   archivedAt: string;
 components: MealComponent[];
 createdAt: string;
 createdByUser: string;
 description: string;
 eligibleForMealPlans: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 id: string;
 lastUpdatedAt: string;
 name: string;
constructor(input: Partial<Meal> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.components = input.components || [];
 this.createdAt = input.createdAt || '';
 this.createdByUser = input.createdByUser || '';
 this.description = input.description || '';
 this.eligibleForMealPlans = input.eligibleForMealPlans || false;
 this.estimatedPortions = input.estimatedPortions || { min: 0 };
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.name = input.name || '';
}
}