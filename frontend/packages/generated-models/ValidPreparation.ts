// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparation {
  lastUpdatedAt?: string;
  restrictToIngredients: boolean;
  id: string;
  instrumentCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  pastTense: string;
  yieldsNothing: boolean;
  conditionExpressionRequired: boolean;
  createdAt: string;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  slug: string;
  vesselCount: NumberRangeWithOptionalMax;
  archivedAt?: string;
  consumesVessel: boolean;
  name: string;
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
}

export class ValidPreparation implements IValidPreparation {
  lastUpdatedAt?: string;
  restrictToIngredients: boolean;
  id: string;
  instrumentCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  pastTense: string;
  yieldsNothing: boolean;
  conditionExpressionRequired: boolean;
  createdAt: string;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  slug: string;
  vesselCount: NumberRangeWithOptionalMax;
  archivedAt?: string;
  consumesVessel: boolean;
  name: string;
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  constructor(input: Partial<ValidPreparation> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.id = input.id = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.onlyForVessels = input.onlyForVessels = false;
    this.pastTense = input.pastTense = '';
    this.yieldsNothing = input.yieldsNothing = false;
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.slug = input.slug = '';
    this.vesselCount = input.vesselCount = { min: 0 };
    this.archivedAt = input.archivedAt;
    this.consumesVessel = input.consumesVessel = false;
    this.name = input.name = '';
    this.temperatureRequired = input.temperatureRequired = false;
    this.timeEstimateRequired = input.timeEstimateRequired = false;
  }
}
