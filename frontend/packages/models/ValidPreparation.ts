// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparation {
  yieldsNothing: boolean;
  description: string;
  ingredientCount: NumberRangeWithOptionalMax;
  temperatureRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  archivedAt: string;
  onlyForVessels: boolean;
  restrictToIngredients: boolean;
  slug: string;
  name: string;
  pastTense: string;
  timeEstimateRequired: boolean;
  consumesVessel: boolean;
  createdAt: string;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  conditionExpressionRequired: boolean;
  iconPath: string;
  id: string;
}

export class ValidPreparation implements IValidPreparation {
  yieldsNothing: boolean;
  description: string;
  ingredientCount: NumberRangeWithOptionalMax;
  temperatureRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  archivedAt: string;
  onlyForVessels: boolean;
  restrictToIngredients: boolean;
  slug: string;
  name: string;
  pastTense: string;
  timeEstimateRequired: boolean;
  consumesVessel: boolean;
  createdAt: string;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  conditionExpressionRequired: boolean;
  iconPath: string;
  id: string;
  constructor(input: Partial<ValidPreparation> = {}) {
    this.yieldsNothing = input.yieldsNothing || false;
    this.description = input.description || '';
    this.ingredientCount = input.ingredientCount || { min: 0 };
    this.temperatureRequired = input.temperatureRequired || false;
    this.vesselCount = input.vesselCount || { min: 0 };
    this.archivedAt = input.archivedAt || '';
    this.onlyForVessels = input.onlyForVessels || false;
    this.restrictToIngredients = input.restrictToIngredients || false;
    this.slug = input.slug || '';
    this.name = input.name || '';
    this.pastTense = input.pastTense || '';
    this.timeEstimateRequired = input.timeEstimateRequired || false;
    this.consumesVessel = input.consumesVessel || false;
    this.createdAt = input.createdAt || '';
    this.instrumentCount = input.instrumentCount || { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.conditionExpressionRequired = input.conditionExpressionRequired || false;
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
  }
}
