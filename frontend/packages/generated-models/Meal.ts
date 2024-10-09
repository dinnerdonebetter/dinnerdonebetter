// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponent } from './MealComponent';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMeal {
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  components: MealComponent;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  createdAt: string;
  createdByUser: string;
  description: string;
}

export class Meal implements IMeal {
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  components: MealComponent;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  createdAt: string;
  createdByUser: string;
  description: string;
  constructor(input: Partial<Meal> = {}) {
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.components = input.components = new MealComponent();
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.name = input.name = '';
    this.createdAt = input.createdAt = '';
    this.createdByUser = input.createdByUser = '';
    this.description = input.description = '';
  }
}
