// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';

export interface IMealPlanEvent {
  options: MealPlanOption;
  createdAt: string;
  endsAt: string;
  id: string;
  mealName: string;
  startsAt: string;
  archivedAt?: string;
  belongsToMealPlan: string;
  lastUpdatedAt?: string;
  notes: string;
}

export class MealPlanEvent implements IMealPlanEvent {
  options: MealPlanOption;
  createdAt: string;
  endsAt: string;
  id: string;
  mealName: string;
  startsAt: string;
  archivedAt?: string;
  belongsToMealPlan: string;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<MealPlanEvent> = {}) {
    this.options = input.options = new MealPlanOption();
    this.createdAt = input.createdAt = '';
    this.endsAt = input.endsAt = '';
    this.id = input.id = '';
    this.mealName = input.mealName = '';
    this.startsAt = input.startsAt = '';
    this.archivedAt = input.archivedAt;
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
