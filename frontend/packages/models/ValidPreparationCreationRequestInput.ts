// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparationCreationRequestInput {
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  instrumentCount: NumberRangeWithOptionalMax;
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

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  instrumentCount: NumberRangeWithOptionalMax;
  name: string;
  onlyForVessels: boolean;
  pastTense: string;
  restrictToIngredients: boolean;
  slug: string;
  temperatureRequired: boolean;
  timeEstimateRequired: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
    this.conditionExpressionRequired = input.conditionExpressionRequired || false;
    this.consumesVessel = input.consumesVessel || false;
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.ingredientCount = input.ingredientCount || { min: 0 };
    this.instrumentCount = input.instrumentCount || { min: 0 };
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
