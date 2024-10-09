// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal';
import { MealPlanOptionVote } from './MealPlanOptionVote';

export interface IMealPlanOption {
  assignedDishwasher?: string;
  id: string;
  lastUpdatedAt?: string;
  meal: Meal;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  archivedAt?: string;
  assignedCook?: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  createdAt: string;
  mealScale: number;
}

export class MealPlanOption implements IMealPlanOption {
  assignedDishwasher?: string;
  id: string;
  lastUpdatedAt?: string;
  meal: Meal;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote;
  archivedAt?: string;
  assignedCook?: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  createdAt: string;
  mealScale: number;
  constructor(input: Partial<MealPlanOption> = {}) {
    this.assignedDishwasher = input.assignedDishwasher;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.meal = input.meal = new Meal();
    this.notes = input.notes = '';
    this.tieBroken = input.tieBroken = false;
    this.votes = input.votes = new MealPlanOptionVote();
    this.archivedAt = input.archivedAt;
    this.assignedCook = input.assignedCook;
    this.belongsToMealPlanEvent = input.belongsToMealPlanEvent = '';
    this.chosen = input.chosen = false;
    this.createdAt = input.createdAt = '';
    this.mealScale = input.mealScale = 0;
  }
}
