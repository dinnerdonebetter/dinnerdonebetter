// GENERATED CODE, DO NOT EDIT MANUALLY



export interface ICreateMealPlanTasksResponse {
   success: boolean;

}

export class CreateMealPlanTasksResponse implements ICreateMealPlanTasksResponse {
   success: boolean;
constructor(input: Partial<CreateMealPlanTasksResponse> = {}) {
	 this.success = input.success || false;
}
}