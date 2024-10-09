// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanOption } from './MealPlanOption';


export interface IMealPlanEvent {
   createdAt: string;
 endsAt: string;
 options: MealPlanOption;
 belongsToMealPlan: string;
 id: string;
 lastUpdatedAt?: string;
 mealName: string;
 notes: string;
 startsAt: string;
 archivedAt?: string;

}

export class MealPlanEvent implements IMealPlanEvent {
   createdAt: string;
 endsAt: string;
 options: MealPlanOption;
 belongsToMealPlan: string;
 id: string;
 lastUpdatedAt?: string;
 mealName: string;
 notes: string;
 startsAt: string;
 archivedAt?: string;
constructor(input: Partial<MealPlanEvent> = {}) {
	 this.createdAt = input.createdAt = '';
 this.endsAt = input.endsAt = '';
 this.options = input.options = new MealPlanOption();
 this.belongsToMealPlan = input.belongsToMealPlan = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.mealName = input.mealName = '';
 this.notes = input.notes = '';
 this.startsAt = input.startsAt = '';
 this.archivedAt = input.archivedAt;
}
}