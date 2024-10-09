// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';

export interface IMealPlanEvent {
  mealName: string;
  notes: string;
  startsAt: string;
  createdAt: string;
  id: string;
  endsAt: string;
  lastUpdatedAt?: string;
  options: MealPlanOption;
  archivedAt?: string;
  belongsToMealPlan: string;
}

export class MealPlanEvent implements IMealPlanEvent {
  mealName: string;
  notes: string;
  startsAt: string;
  createdAt: string;
  id: string;
  endsAt: string;
  lastUpdatedAt?: string;
  options: MealPlanOption;
  archivedAt?: string;
  belongsToMealPlan: string;
  constructor(input: Partial<MealPlanEvent> = {}) {
    this.mealName = input.mealName = '';
    this.notes = input.notes = '';
    this.startsAt = input.startsAt = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.endsAt = input.endsAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.options = input.options = new MealPlanOption();
    this.archivedAt = input.archivedAt;
    this.belongsToMealPlan = input.belongsToMealPlan = '';
  }
}
