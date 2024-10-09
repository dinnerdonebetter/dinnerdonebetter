// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEvent } from './MealPlanEvent';
import { ValidMealPlanStatus, ValidMealPlanElectionMethod } from './enums';

export interface IMealPlan {
  archivedAt?: string;
  belongsToHousehold: string;
  createdBy: string;
  lastUpdatedAt?: string;
  notes: string;
  status: ValidMealPlanStatus;
  votingDeadline: string;
  createdAt: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent;
  groceryListInitialized: boolean;
  id: string;
  tasksCreated: boolean;
}

export class MealPlan implements IMealPlan {
  archivedAt?: string;
  belongsToHousehold: string;
  createdBy: string;
  lastUpdatedAt?: string;
  notes: string;
  status: ValidMealPlanStatus;
  votingDeadline: string;
  createdAt: string;
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEvent;
  groceryListInitialized: boolean;
  id: string;
  tasksCreated: boolean;
  constructor(input: Partial<MealPlan> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.createdBy = input.createdBy = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.status = input.status = 'awaiting_votes';
    this.votingDeadline = input.votingDeadline = '';
    this.createdAt = input.createdAt = '';
    this.electionMethod = input.electionMethod = 'schulze';
    this.events = input.events = new MealPlanEvent();
    this.groceryListInitialized = input.groceryListInitialized = false;
    this.id = input.id = '';
    this.tasksCreated = input.tasksCreated = false;
  }
}
