// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IInitializeMealPlanGroceryListRequest {
   householdID: string;

}

export class InitializeMealPlanGroceryListRequest implements IInitializeMealPlanGroceryListRequest {
   householdID: string;
constructor(input: Partial<InitializeMealPlanGroceryListRequest> = {}) {
	 this.householdID = input.householdID || '';
}
}