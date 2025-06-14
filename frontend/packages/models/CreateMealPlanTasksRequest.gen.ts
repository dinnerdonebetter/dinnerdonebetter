// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ICreateMealPlanTasksRequest {
  accountID: string;
}

export class CreateMealPlanTasksRequest implements ICreateMealPlanTasksRequest {
  accountID: string;
  constructor(input: Partial<CreateMealPlanTasksRequest> = {}) {
    this.accountID = input.accountID || '';
  }
}
