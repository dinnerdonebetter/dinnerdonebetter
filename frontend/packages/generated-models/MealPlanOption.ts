// GENERATED CODE, DO NOT EDIT MANUALLY

 import { Meal } from './Meal';
 import { MealPlanOptionVote } from './MealPlanOptionVote';


export interface IMealPlanOption {
   createdAt: string;
 id: string;
 mealScale: number;
 tieBroken: boolean;
 assignedCook?: string;
 assignedDishwasher?: string;
 chosen: boolean;
 meal: Meal;
 notes: string;
 votes: MealPlanOptionVote;
 archivedAt?: string;
 belongsToMealPlanEvent: string;
 lastUpdatedAt?: string;

}

export class MealPlanOption implements IMealPlanOption {
   createdAt: string;
 id: string;
 mealScale: number;
 tieBroken: boolean;
 assignedCook?: string;
 assignedDishwasher?: string;
 chosen: boolean;
 meal: Meal;
 notes: string;
 votes: MealPlanOptionVote;
 archivedAt?: string;
 belongsToMealPlanEvent: string;
 lastUpdatedAt?: string;
constructor(input: Partial<MealPlanOption> = {}) {
	 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.mealScale = input.mealScale = 0;
 this.tieBroken = input.tieBroken = false;
 this.assignedCook = input.assignedCook;
 this.assignedDishwasher = input.assignedDishwasher;
 this.chosen = input.chosen = false;
 this.meal = input.meal = new Meal();
 this.notes = input.notes = '';
 this.votes = input.votes = new MealPlanOptionVote();
 this.archivedAt = input.archivedAt;
 this.belongsToMealPlanEvent = input.belongsToMealPlanEvent = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
}
}