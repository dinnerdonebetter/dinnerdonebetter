// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEvent } from './MealPlanEvent';
import { ValidMealPlanStatus, ValidMealPlanElectionMethod } from './enums';

export interface IMealPlan {
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  votingDeadline: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent;
  lastUpdatedAt?: string;
  notes: string;
  groceryListInitialized: boolean;
  id: string;
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
  createdBy: string;
}

export class MealPlan implements IMealPlan {
  status: ValidMealPlanStatus;
  tasksCreated: boolean;
  votingDeadline: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent;
  lastUpdatedAt?: string;
  notes: string;
  groceryListInitialized: boolean;
  id: string;
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
  createdBy: string;
  constructor(input: Partial<MealPlan> = {}) {
    this.status = input.status = 'awaiting_votes';
    this.tasksCreated = input.tasksCreated = false;
    this.votingDeadline = input.votingDeadline = '';
    this.electionMethod = input.electionMethod = 'schulze';
    this.events = input.events = new MealPlanEvent();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.groceryListInitialized = input.groceryListInitialized = false;
    this.id = input.id = '';
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.createdAt = input.createdAt = '';
    this.createdBy = input.createdBy = '';
  }
}
