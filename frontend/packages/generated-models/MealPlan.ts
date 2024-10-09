// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanEvent } from './MealPlanEvent';
 import { ValidMealPlanStatus, ValidMealPlanElectionMethod } from './enums';


export interface IMealPlan {
   archivedAt?: string;
 belongsToHousehold: string;
 events: MealPlanEvent;
 notes: string;
 status: ValidMealPlanStatus;
 tasksCreated: boolean;
 createdAt: string;
 createdBy: string;
 electionMethod: ValidMealPlanElectionMethod;
 groceryListInitialized: boolean;
 id: string;
 lastUpdatedAt?: string;
 votingDeadline: string;

}

export class MealPlan implements IMealPlan {
   archivedAt?: string;
 belongsToHousehold: string;
 events: MealPlanEvent;
 notes: string;
 status: ValidMealPlanStatus;
 tasksCreated: boolean;
 createdAt: string;
 createdBy: string;
 electionMethod: ValidMealPlanElectionMethod;
 groceryListInitialized: boolean;
 id: string;
 lastUpdatedAt?: string;
 votingDeadline: string;
constructor(input: Partial<MealPlan> = {}) {
	 this.archivedAt = input.archivedAt;
 this.belongsToHousehold = input.belongsToHousehold = '';
 this.events = input.events = new MealPlanEvent();
 this.notes = input.notes = '';
 this.status = input.status = 'awaiting_votes';
 this.tasksCreated = input.tasksCreated = false;
 this.createdAt = input.createdAt = '';
 this.createdBy = input.createdBy = '';
 this.electionMethod = input.electionMethod = 'schulze';
 this.groceryListInitialized = input.groceryListInitialized = false;
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.votingDeadline = input.votingDeadline = '';
}
}