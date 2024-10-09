// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ICreateMealPlanTasksRequest {
  householdID: string;
}

export class CreateMealPlanTasksRequest implements ICreateMealPlanTasksRequest {
  householdID: string;
  constructor(input: Partial<CreateMealPlanTasksRequest> = {}) {
    this.householdID = input.householdID = '';
  }
}
