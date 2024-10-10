// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparation {
  archivedAt: string;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  createdAt: string;
  description: string;
  iconPath: string;
  id: string;
  ingredientCount: NumberRangeWithOptionalMax;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  name: string;
  onlyForVessels: boolean;
  pastTense: string;
  restrictToIngredients: boolean;
  slug: string;
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
}

export class ValidPreparation implements IValidPreparation {
  archivedAt: string;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  createdAt: string;
  description: string;
  iconPath: string;
  id: string;
  ingredientCount: NumberRangeWithOptionalMax;
  instrumentCount: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  name: string;
  onlyForVessels: boolean;
  pastTense: string;
  restrictToIngredients: boolean;
  slug: string;
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  constructor(input: Partial<ValidPreparation> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.conditionExpressionRequired = input.conditionExpressionRequired || false;
    this.consumesVessel = input.consumesVessel || false;
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
    this.ingredientCount = input.ingredientCount || { min: 0 };
    this.instrumentCount = input.instrumentCount || { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.onlyForVessels = input.onlyForVessels || false;
    this.pastTense = input.pastTense || '';
    this.restrictToIngredients = input.restrictToIngredients || false;
    this.slug = input.slug || '';
    this.temperatureRequired = input.temperatureRequired || false;
    this.timeEstimateRequired = input.timeEstimateRequired || false;
    this.vesselCount = input.vesselCount || { min: 0 };
    this.yieldsNothing = input.yieldsNothing || false;
  }
}
