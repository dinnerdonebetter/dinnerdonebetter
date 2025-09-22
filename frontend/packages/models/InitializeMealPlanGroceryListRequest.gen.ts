// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IInitializeMealPlanGroceryListRequest {
  accountID: string;
}

export class InitializeMealPlanGroceryListRequest implements IInitializeMealPlanGroceryListRequest {
  accountID: string;
  constructor(input: Partial<InitializeMealPlanGroceryListRequest> = {}) {
    this.accountID = input.accountID || '';
  }
}
