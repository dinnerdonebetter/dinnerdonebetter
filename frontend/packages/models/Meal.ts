// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponent } from './MealComponent';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMeal {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  components: MealComponent[];
  description: string;
  eligibleForMealPlans: boolean;
  id: string;
  name: string;
}

export class Meal implements IMeal {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  components: MealComponent[];
  description: string;
  eligibleForMealPlans: boolean;
  id: string;
  name: string;
  constructor(input: Partial<Meal> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.createdByUser = input.createdByUser || '';
    this.estimatedPortions = input.estimatedPortions || { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.components = input.components || [];
    this.description = input.description || '';
    this.eligibleForMealPlans = input.eligibleForMealPlans || false;
    this.id = input.id || '';
    this.name = input.name || '';
  }
}
