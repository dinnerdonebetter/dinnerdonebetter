// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparationCreationRequestInput {
  description: string;
  vesselCount: NumberRangeWithOptionalMax;
  conditionExpressionRequired: boolean;
  ingredientCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  pastTense: string;
  timeEstimateRequired: boolean;
  yieldsNothing: boolean;
  restrictToIngredients: boolean;
  temperatureRequired: boolean;
  consumesVessel: boolean;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  name: string;
  slug: string;
}

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
  description: string;
  vesselCount: NumberRangeWithOptionalMax;
  conditionExpressionRequired: boolean;
  ingredientCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  pastTense: string;
  timeEstimateRequired: boolean;
  yieldsNothing: boolean;
  restrictToIngredients: boolean;
  temperatureRequired: boolean;
  consumesVessel: boolean;
  iconPath: string;
  instrumentCount: NumberRangeWithOptionalMax;
  name: string;
  slug: string;
  constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
    this.description = input.description || '';
    this.vesselCount = input.vesselCount || { min: 0 };
    this.conditionExpressionRequired = input.conditionExpressionRequired || false;
    this.ingredientCount = input.ingredientCount || { min: 0 };
    this.onlyForVessels = input.onlyForVessels || false;
    this.pastTense = input.pastTense || '';
    this.timeEstimateRequired = input.timeEstimateRequired || false;
    this.yieldsNothing = input.yieldsNothing || false;
    this.restrictToIngredients = input.restrictToIngredients || false;
    this.temperatureRequired = input.temperatureRequired || false;
    this.consumesVessel = input.consumesVessel || false;
    this.iconPath = input.iconPath || '';
    this.instrumentCount = input.instrumentCount || { min: 0 };
    this.name = input.name || '';
    this.slug = input.slug || '';
  }
}
