// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparation {
  restrictToIngredients: boolean;
  temperatureRequired: boolean;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  id: string;
  onlyForVessels: boolean;
  pastTense: string;
  archivedAt?: string;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  slug: string;
  description: string;
  ingredientCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  createdAt: string;
  name: string;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
}

export class ValidPreparation implements IValidPreparation {
  restrictToIngredients: boolean;
  temperatureRequired: boolean;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  id: string;
  onlyForVessels: boolean;
  pastTense: string;
  archivedAt?: string;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  slug: string;
  description: string;
  ingredientCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  createdAt: string;
  name: string;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  constructor(input: Partial<ValidPreparation> = {}) {
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.temperatureRequired = input.temperatureRequired = false;
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.consumesVessel = input.consumesVessel = false;
    this.id = input.id = '';
    this.onlyForVessels = input.onlyForVessels = false;
    this.pastTense = input.pastTense = '';
    this.archivedAt = input.archivedAt;
    this.iconPath = input.iconPath = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.slug = input.slug = '';
    this.description = input.description = '';
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.yieldsNothing = input.yieldsNothing = false;
    this.createdAt = input.createdAt = '';
    this.name = input.name = '';
    this.timeEstimateRequired = input.timeEstimateRequired = false;
    this.vesselCount = input.vesselCount = { min: 0 };
  }
}
