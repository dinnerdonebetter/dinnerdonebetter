// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal';
import { MealPlanOptionVote } from './MealPlanOptionVote';

export interface IMealPlanOption {
  id: string;
  lastUpdatedAt: string;
  mealScale: number;
  tieBroken: boolean;
  assignedDishwasher: string;
  chosen: boolean;
  belongsToMealPlanEvent: string;
  createdAt: string;
  meal: Meal;
  notes: string;
  votes: MealPlanOptionVote[];
  archivedAt: string;
  assignedCook: string;
}

export class MealPlanOption implements IMealPlanOption {
  id: string;
  lastUpdatedAt: string;
  mealScale: number;
  tieBroken: boolean;
  assignedDishwasher: string;
  chosen: boolean;
  belongsToMealPlanEvent: string;
  createdAt: string;
  meal: Meal;
  notes: string;
  votes: MealPlanOptionVote[];
  archivedAt: string;
  assignedCook: string;
  constructor(input: Partial<MealPlanOption> = {}) {
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mealScale = input.mealScale || 0;
    this.tieBroken = input.tieBroken || false;
    this.assignedDishwasher = input.assignedDishwasher || '';
    this.chosen = input.chosen || false;
    this.belongsToMealPlanEvent = input.belongsToMealPlanEvent || '';
    this.createdAt = input.createdAt || '';
    this.meal = input.meal || new Meal();
    this.notes = input.notes || '';
    this.votes = input.votes || [];
    this.archivedAt = input.archivedAt || '';
    this.assignedCook = input.assignedCook || '';
  }
}
