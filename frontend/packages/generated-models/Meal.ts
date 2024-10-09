// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponent } from './MealComponent';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMeal {
  createdByUser: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  archivedAt?: string;
  components: MealComponent;
  createdAt: string;
  eligibleForMealPlans: boolean;
  lastUpdatedAt?: string;
  name: string;
}

export class Meal implements IMeal {
  createdByUser: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  archivedAt?: string;
  components: MealComponent;
  createdAt: string;
  eligibleForMealPlans: boolean;
  lastUpdatedAt?: string;
  name: string;
  constructor(input: Partial<Meal> = {}) {
    this.createdByUser = input.createdByUser = '';
    this.description = input.description = '';
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.id = input.id = '';
    this.archivedAt = input.archivedAt;
    this.components = input.components = new MealComponent();
    this.createdAt = input.createdAt = '';
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
  }
}
