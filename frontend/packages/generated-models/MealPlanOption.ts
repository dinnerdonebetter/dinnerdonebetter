// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal';
import { MealPlanOptionVote } from './MealPlanOptionVote';

export interface IMealPlanOption {
  lastUpdatedAt?: string;
  mealScale: number;
  assignedCook?: string;
  assignedDishwasher?: string;
  createdAt: string;
  id: string;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  archivedAt?: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  meal: Meal;
}

export class MealPlanOption implements IMealPlanOption {
  lastUpdatedAt?: string;
  mealScale: number;
  assignedCook?: string;
  assignedDishwasher?: string;
  createdAt: string;
  id: string;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  archivedAt?: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  meal: Meal;
  constructor(input: Partial<MealPlanOption> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.mealScale = input.mealScale = 0;
    this.assignedCook = input.assignedCook;
    this.assignedDishwasher = input.assignedDishwasher;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.tieBroken = input.tieBroken = false;
    this.votes = input.votes = new MealPlanOptionVote();
    this.archivedAt = input.archivedAt;
    this.belongsToMealPlanEvent = input.belongsToMealPlanEvent = '';
    this.chosen = input.chosen = false;
    this.meal = input.meal = new Meal();
  }
}
