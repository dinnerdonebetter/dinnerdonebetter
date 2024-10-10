// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal';
import { MealPlanOptionVote } from './MealPlanOptionVote';

export interface IMealPlanOption {
  archivedAt: string;
  assignedCook: string;
  assignedDishwasher: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  meal: Meal;
  mealScale: number;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote[];
}

export class MealPlanOption implements IMealPlanOption {
  archivedAt: string;
  assignedCook: string;
  assignedDishwasher: string;
  belongsToMealPlanEvent: string;
  chosen: boolean;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  meal: Meal;
  mealScale: number;
  notes: string;
  tieBroken: boolean;
  votes: MealPlanOptionVote[];
  constructor(input: Partial<MealPlanOption> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.assignedCook = input.assignedCook || '';
    this.assignedDishwasher = input.assignedDishwasher || '';
    this.belongsToMealPlanEvent = input.belongsToMealPlanEvent || '';
    this.chosen = input.chosen || false;
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.meal = input.meal || new Meal();
    this.mealScale = input.mealScale || 0;
    this.notes = input.notes || '';
    this.tieBroken = input.tieBroken || false;
    this.votes = input.votes || [];
  }
}
