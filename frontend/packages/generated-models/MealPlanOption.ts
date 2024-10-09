// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal';
import { MealPlanOptionVote } from './MealPlanOptionVote';

export interface IMealPlanOption {
  meal: Meal;
  mealScale: number;
  archivedAt?: string;
  createdAt: string;
  id: string;
  chosen: boolean;
  lastUpdatedAt?: string;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  assignedCook?: string;
  assignedDishwasher?: string;
  belongsToMealPlanEvent: string;
}

export class MealPlanOption implements IMealPlanOption {
  meal: Meal;
  mealScale: number;
  archivedAt?: string;
  createdAt: string;
  id: string;
  chosen: boolean;
  lastUpdatedAt?: string;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  assignedCook?: string;
  assignedDishwasher?: string;
  belongsToMealPlanEvent: string;
  constructor(input: Partial<MealPlanOption> = {}) {
    this.meal = input.meal = new Meal();
    this.mealScale = input.mealScale = 0;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.chosen = input.chosen = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.tieBroken = input.tieBroken = false;
    this.votes = input.votes = new MealPlanOptionVote();
    this.assignedCook = input.assignedCook;
    this.assignedDishwasher = input.assignedDishwasher;
    this.belongsToMealPlanEvent = input.belongsToMealPlanEvent = '';
  }
}
