// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanOptionCreationRequestInput } from './MealPlanOptionCreationRequestInput';


export interface IMealPlanEventCreationRequestInput {
   mealName: string;
 notes: string;
 options: MealPlanOptionCreationRequestInput;
 startsAt: string;
 endsAt: string;

}

export class MealPlanEventCreationRequestInput implements IMealPlanEventCreationRequestInput {
   mealName: string;
 notes: string;
 options: MealPlanOptionCreationRequestInput;
 startsAt: string;
 endsAt: string;
constructor(input: Partial<MealPlanEventCreationRequestInput> = {}) {
	 this.mealName = input.mealName = '';
 this.notes = input.notes = '';
 this.options = input.options = new MealPlanOptionCreationRequestInput();
 this.startsAt = input.startsAt = '';
 this.endsAt = input.endsAt = '';
}
}