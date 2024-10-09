// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponent } from './MealComponent';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMeal {
  archivedAt?: string;
  createdAt: string;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  components: MealComponent;
  createdByUser: string;
  description: string;
  id: string;
  name: string;
}

export class Meal implements IMeal {
  archivedAt?: string;
  createdAt: string;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  components: MealComponent;
  createdByUser: string;
  description: string;
  id: string;
  name: string;
  constructor(input: Partial<Meal> = {}) {
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.components = input.components = new MealComponent();
    this.createdByUser = input.createdByUser = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.name = input.name = '';
  }
}
