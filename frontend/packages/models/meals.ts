// Code generated by gen_typescript. DO NOT EDIT.

import { NumberRangeWithOptionalMax, OptionalNumberRange } from './main';
import { MealComponent, MealComponentCreationRequestInput, MealComponentUpdateRequestInput } from './mealComponents';

export interface IMeal {
  createdAt: NonNullable<string>;
  archivedAt?: string;
  lastUpdatedAt?: string;
  estimatedPortions: NonNullable<NumberRangeWithOptionalMax>;
  id: NonNullable<string>;
  description: NonNullable<string>;
  createdByUser: NonNullable<string>;
  name: NonNullable<string>;
  components: NonNullable<Array<MealComponent>>;
  eligibleForMealPlans: NonNullable<boolean>;
}

export class Meal implements IMeal {
  createdAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  archivedAt?: string;
  lastUpdatedAt?: string;
  estimatedPortions: NonNullable<NumberRangeWithOptionalMax> = { min: 0 };
  id: NonNullable<string> = '';
  description: NonNullable<string> = '';
  createdByUser: NonNullable<string> = '';
  name: NonNullable<string> = '';
  components: NonNullable<Array<MealComponent>> = [];
  eligibleForMealPlans: NonNullable<boolean> = false;

  constructor(input: Partial<Meal> = {}) {
    this.createdAt = input.createdAt ?? '1970-01-01T00:00:00Z';
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.estimatedPortions = input.estimatedPortions ?? { min: 0 };
    this.id = input.id ?? '';
    this.description = input.description ?? '';
    this.createdByUser = input.createdByUser ?? '';
    this.name = input.name ?? '';
    this.components = input.components ?? [];
    this.eligibleForMealPlans = input.eligibleForMealPlans ?? false;
  }
}

export interface IMealCreationRequestInput {
  estimatedPortions: NonNullable<NumberRangeWithOptionalMax>;
  name: NonNullable<string>;
  description: NonNullable<string>;
  components: NonNullable<Array<MealComponentCreationRequestInput>>;
  eligibleForMealPlans: NonNullable<boolean>;
}

export class MealCreationRequestInput implements IMealCreationRequestInput {
  estimatedPortions: NonNullable<NumberRangeWithOptionalMax> = { min: 0 };
  name: NonNullable<string> = '';
  description: NonNullable<string> = '';
  components: NonNullable<Array<MealComponentCreationRequestInput>> = [];
  eligibleForMealPlans: NonNullable<boolean> = false;

  constructor(input: Partial<MealCreationRequestInput> = {}) {
    this.estimatedPortions = input.estimatedPortions ?? { min: 0 };
    this.name = input.name ?? '';
    this.description = input.description ?? '';
    this.components = input.components ?? [];
    this.eligibleForMealPlans = input.eligibleForMealPlans ?? false;
  }
}

export interface IMealUpdateRequestInput {
  name?: string;
  description?: string;
  estimatedPortions: NonNullable<OptionalNumberRange>;
  eligibleForMealPlans?: boolean;
  recipes: NonNullable<Array<MealComponentUpdateRequestInput>>;
}

export class MealUpdateRequestInput implements IMealUpdateRequestInput {
  name?: string;
  description?: string;
  estimatedPortions: NonNullable<OptionalNumberRange> = {};
  eligibleForMealPlans?: boolean = false;
  recipes: NonNullable<Array<MealComponentUpdateRequestInput>> = [];

  constructor(input: Partial<MealUpdateRequestInput> = {}) {
    this.name = input.name;
    this.description = input.description;
    this.estimatedPortions = input.estimatedPortions ?? {};
    this.eligibleForMealPlans = input.eligibleForMealPlans ?? false;
    this.recipes = input.recipes ?? [];
  }
}
