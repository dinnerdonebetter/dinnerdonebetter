// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparationCreationRequestInput {
  vesselCount: NumberRangeWithOptionalMax;
  slug: string;
  instrumentCount: NumberRangeWithOptionalMax;
  name: string;
  onlyForVessels: boolean;
  pastTense: string;
  timeEstimateRequired: boolean;
  conditionExpressionRequired: boolean;
  temperatureRequired: boolean;
  restrictToIngredients: boolean;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  consumesVessel: boolean;
}

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
  vesselCount: NumberRangeWithOptionalMax;
  slug: string;
  instrumentCount: NumberRangeWithOptionalMax;
  name: string;
  onlyForVessels: boolean;
  pastTense: string;
  timeEstimateRequired: boolean;
  conditionExpressionRequired: boolean;
  temperatureRequired: boolean;
  restrictToIngredients: boolean;
  description: string;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  yieldsNothing: boolean;
  consumesVessel: boolean;
  constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
    this.vesselCount = input.vesselCount = { min: 0 };
    this.slug = input.slug = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.name = input.name = '';
    this.onlyForVessels = input.onlyForVessels = false;
    this.pastTense = input.pastTense = '';
    this.timeEstimateRequired = input.timeEstimateRequired = false;
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.temperatureRequired = input.temperatureRequired = false;
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.yieldsNothing = input.yieldsNothing = false;
    this.consumesVessel = input.consumesVessel = false;
  }
}
