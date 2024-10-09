// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEvent } from './MealPlanEvent';
import { ValidMealPlanStatus, ValidMealPlanElectionMethod } from './enums';

export interface IMealPlan {
  id: string;
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  archivedAt: string;
  createdBy: string;
  events: MealPlanEvent[];
  groceryListInitialized: boolean;
  lastUpdatedAt: string;
  notes: string;
  votingDeadline: string;
  belongsToHousehold: string;
  createdAt: string;
  electionMethod: ValidMealPlanElectionMethod;
}

export class MealPlan implements IMealPlan {
  id: string;
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  archivedAt: string;
  createdBy: string;
  events: MealPlanEvent[];
  groceryListInitialized: boolean;
  lastUpdatedAt: string;
  notes: string;
  votingDeadline: string;
  belongsToHousehold: string;
  createdAt: string;
  electionMethod: ValidMealPlanElectionMethod;
  constructor(input: Partial<MealPlan> = {}) {
    this.id = input.id || '';
    this.status = input.status || 'awaiting_votes';
    this.tasksCreated = input.tasksCreated || false;
    this.archivedAt = input.archivedAt || '';
    this.createdBy = input.createdBy || '';
    this.events = input.events || [];
    this.groceryListInitialized = input.groceryListInitialized || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.votingDeadline = input.votingDeadline || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.createdAt = input.createdAt || '';
    this.electionMethod = input.electionMethod || 'schulze';
  }
}
