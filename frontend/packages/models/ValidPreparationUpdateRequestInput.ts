// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidPreparationUpdateRequestInput {
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  instrumentCount: OptionalNumberRange;
  restrictToIngredients: boolean;
  slug: string;
  consumesVessel: boolean;
  ingredientCount: OptionalNumberRange;
  name: string;
  vesselCount: OptionalNumberRange;
  yieldsNothing: boolean;
  description: string;
  iconPath: string;
  onlyForVessels: boolean;
  pastTense: string;
  conditionExpressionRequired: boolean;
}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  instrumentCount: OptionalNumberRange;
  restrictToIngredients: boolean;
  slug: string;
  consumesVessel: boolean;
  ingredientCount: OptionalNumberRange;
  name: string;
  vesselCount: OptionalNumberRange;
  yieldsNothing: boolean;
  description: string;
  iconPath: string;
  onlyForVessels: boolean;
  pastTense: string;
  conditionExpressionRequired: boolean;
  constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
    this.temperatureRequired = input.temperatureRequired || false;
    this.timeEstimateRequired = input.timeEstimateRequired || false;
    this.instrumentCount = input.instrumentCount || {};
    this.restrictToIngredients = input.restrictToIngredients || false;
    this.slug = input.slug || '';
    this.consumesVessel = input.consumesVessel || false;
    this.ingredientCount = input.ingredientCount || {};
    this.name = input.name || '';
    this.vesselCount = input.vesselCount || {};
    this.yieldsNothing = input.yieldsNothing || false;
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.onlyForVessels = input.onlyForVessels || false;
    this.pastTense = input.pastTense || '';
    this.conditionExpressionRequired = input.conditionExpressionRequired || false;
  }
}
