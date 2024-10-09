// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanOption } from './MealPlanOption';
 import { RecipePrepTask } from './RecipePrepTask';
 import { MealPlanTaskStatus } from './enums';


export interface IMealPlanTask {
   assignedToUser?: string;
 completedAt?: string;
 lastUpdatedAt?: string;
 recipePrepTask: RecipePrepTask;
 status: MealPlanTaskStatus;
 createdAt: string;
 creationExplanation: string;
 id: string;
 mealPlanOption: MealPlanOption;
 statusExplanation: string;

}

export class MealPlanTask implements IMealPlanTask {
   assignedToUser?: string;
 completedAt?: string;
 lastUpdatedAt?: string;
 recipePrepTask: RecipePrepTask;
 status: MealPlanTaskStatus;
 createdAt: string;
 creationExplanation: string;
 id: string;
 mealPlanOption: MealPlanOption;
 statusExplanation: string;
constructor(input: Partial<MealPlanTask> = {}) {
	 this.assignedToUser = input.assignedToUser;
 this.completedAt = input.completedAt;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.recipePrepTask = input.recipePrepTask = new RecipePrepTask();
 this.status = input.status = 'unfinished';
 this.createdAt = input.createdAt = '';
 this.creationExplanation = input.creationExplanation = '';
 this.id = input.id = '';
 this.mealPlanOption = input.mealPlanOption = new MealPlanOption();
 this.statusExplanation = input.statusExplanation = '';
}
}