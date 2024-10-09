// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparation {
  archivedAt?: string;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  slug: string;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  ingredientCount: NumberRangeWithOptionalMax;
  pastTense: string;
  createdAt: string;
  lastUpdatedAt?: string;
  name: string;
  onlyForVessels: boolean;
  restrictToIngredients: boolean;
  description: string;
  id: string;
  temperatureRequired: boolean;
}

export class ValidPreparation implements IValidPreparation {
  archivedAt?: string;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  slug: string;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  ingredientCount: NumberRangeWithOptionalMax;
  pastTense: string;
  createdAt: string;
  lastUpdatedAt?: string;
  name: string;
  onlyForVessels: boolean;
  restrictToIngredients: boolean;
  description: string;
  id: string;
  temperatureRequired: boolean;
  constructor(input: Partial<ValidPreparation> = {}) {
    this.archivedAt = input.archivedAt;
    this.iconPath = input.iconPath = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.slug = input.slug = '';
    this.timeEstimateRequired = input.timeEstimateRequired = false;
    this.vesselCount = input.vesselCount = { min: 0 };
    this.yieldsNothing = input.yieldsNothing = false;
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.consumesVessel = input.consumesVessel = false;
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.pastTense = input.pastTense = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.onlyForVessels = input.onlyForVessels = false;
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.description = input.description = '';
    this.id = input.id = '';
    this.temperatureRequired = input.temperatureRequired = false;
  }
}
