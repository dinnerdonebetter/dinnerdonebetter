// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';

export interface IMealPlanEvent {
  archivedAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  endsAt: string;
  id: string;
  lastUpdatedAt: string;
  mealName: string;
  notes: string;
  options: MealPlanOption[];
  startsAt: string;
}

export class MealPlanEvent implements IMealPlanEvent {
  archivedAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  endsAt: string;
  id: string;
  lastUpdatedAt: string;
  mealName: string;
  notes: string;
  options: MealPlanOption[];
  startsAt: string;
  constructor(input: Partial<MealPlanEvent> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.createdAt = input.createdAt || '';
    this.endsAt = input.endsAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mealName = input.mealName || '';
    this.notes = input.notes || '';
    this.options = input.options || [];
    this.startsAt = input.startsAt || '';
  }
}
