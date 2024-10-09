// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEvent } from './MealPlanEvent';
import { ValidMealPlanElectionMethod, ValidMealPlanStatus } from './enums';

export interface IMealPlan {
  belongsToHousehold: string;
  createdBy: string;
  electionMethod: ValidMealPlanElectionMethod;
  groceryListInitialized: boolean;
  status: ValidMealPlanStatus;
  votingDeadline: string;
  tasksCreated: boolean;
  archivedAt?: string;
  createdAt: string;
  events: MealPlanEvent;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
}

export class MealPlan implements IMealPlan {
  belongsToHousehold: string;
  createdBy: string;
  electionMethod: ValidMealPlanElectionMethod;
  groceryListInitialized: boolean;
  status: ValidMealPlanStatus;
  votingDeadline: string;
  tasksCreated: boolean;
  archivedAt?: string;
  createdAt: string;
  events: MealPlanEvent;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<MealPlan> = {}) {
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.createdBy = input.createdBy = '';
    this.electionMethod = input.electionMethod = 'schulze';
    this.groceryListInitialized = input.groceryListInitialized = false;
    this.status = input.status = 'awaiting_votes';
    this.votingDeadline = input.votingDeadline = '';
    this.tasksCreated = input.tasksCreated = false;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.events = input.events = new MealPlanEvent();
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
