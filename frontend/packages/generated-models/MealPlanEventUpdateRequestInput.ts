// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IMealPlanEventUpdateRequestInput {
   mealName?: string;
 notes?: string;
 startsAt?: string;
 endsAt?: string;

}

export class MealPlanEventUpdateRequestInput implements IMealPlanEventUpdateRequestInput {
   mealName?: string;
 notes?: string;
 startsAt?: string;
 endsAt?: string;
constructor(input: Partial<MealPlanEventUpdateRequestInput> = {}) {
	 this.mealName = input.mealName;
 this.notes = input.notes;
 this.startsAt = input.startsAt;
 this.endsAt = input.endsAt;
}
}