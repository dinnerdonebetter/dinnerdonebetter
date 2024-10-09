// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';

export interface IMealPlanEvent {
  archivedAt: string;
  endsAt: string;
  lastUpdatedAt: string;
  notes: string;
  startsAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  id: string;
  mealName: string;
  options: MealPlanOption[];
}

export class MealPlanEvent implements IMealPlanEvent {
  archivedAt: string;
  endsAt: string;
  lastUpdatedAt: string;
  notes: string;
  startsAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  id: string;
  mealName: string;
  options: MealPlanOption[];
  constructor(input: Partial<MealPlanEvent> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.endsAt = input.endsAt || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.startsAt = input.startsAt || '';
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.mealName = input.mealName || '';
    this.options = input.options || [];
  }
}
