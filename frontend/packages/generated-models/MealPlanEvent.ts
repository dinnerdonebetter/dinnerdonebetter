// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';

export interface IMealPlanEvent {
  lastUpdatedAt?: string;
  mealName: string;
  notes: string;
  belongsToMealPlan: string;
  endsAt: string;
  id: string;
  startsAt: string;
  archivedAt?: string;
  createdAt: string;
  options: MealPlanOption;
}

export class MealPlanEvent implements IMealPlanEvent {
  lastUpdatedAt?: string;
  mealName: string;
  notes: string;
  belongsToMealPlan: string;
  endsAt: string;
  id: string;
  startsAt: string;
  archivedAt?: string;
  createdAt: string;
  options: MealPlanOption;
  constructor(input: Partial<MealPlanEvent> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.mealName = input.mealName = '';
    this.notes = input.notes = '';
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.endsAt = input.endsAt = '';
    this.id = input.id = '';
    this.startsAt = input.startsAt = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.options = input.options = new MealPlanOption();
  }
}
