// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEvent } from './MealPlanEvent.gen';
import { ValidMealPlanElectionMethod, ValidMealPlanStatus } from './enums.gen';

export interface IMealPlan {
  archivedAt: string;
  belongsToHousehold: string;
  createdAt: string;
  createdBy: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent[];
  groceryListInitialized: boolean;
  id: string;
  lastUpdatedAt: string;
  notes: string;
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  votingDeadline: string;
}

export class MealPlan implements IMealPlan {
  archivedAt: string;
  belongsToHousehold: string;
  createdAt: string;
  createdBy: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent[];
  groceryListInitialized: boolean;
  id: string;
  lastUpdatedAt: string;
  notes: string;
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  votingDeadline: string;
  constructor(input: Partial<MealPlan> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.createdAt = input.createdAt || '';
    this.createdBy = input.createdBy || '';
    this.electionMethod = input.electionMethod || 'schulze';
    this.events = input.events || [];
    this.groceryListInitialized = input.groceryListInitialized || false;
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.status = input.status || 'awaiting_votes';
    this.tasksCreated = input.tasksCreated || false;
    this.votingDeadline = input.votingDeadline || '';
  }
}
