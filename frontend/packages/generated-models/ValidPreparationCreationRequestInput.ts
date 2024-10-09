// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparationCreationRequestInput {
  conditionExpressionRequired: boolean;
  timeEstimateRequired: boolean;
  consumesVessel: boolean;
  description: string;
  instrumentCount: NumberRangeWithOptionalMax;
  pastTense: string;
  restrictToIngredients: boolean;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  vesselCount: NumberRangeWithOptionalMax;
  name: string;
  onlyForVessels: boolean;
  slug: string;
  temperatureRequired: boolean;
  yieldsNothing: boolean;
}

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
  conditionExpressionRequired: boolean;
  timeEstimateRequired: boolean;
  consumesVessel: boolean;
  description: string;
  instrumentCount: NumberRangeWithOptionalMax;
  pastTense: string;
  restrictToIngredients: boolean;
  iconPath: string;
  ingredientCount: NumberRangeWithOptionalMax;
  vesselCount: NumberRangeWithOptionalMax;
  name: string;
  onlyForVessels: boolean;
  slug: string;
  temperatureRequired: boolean;
  yieldsNothing: boolean;
  constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.timeEstimateRequired = input.timeEstimateRequired = false;
    this.consumesVessel = input.consumesVessel = false;
    this.description = input.description = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.pastTense = input.pastTense = '';
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.iconPath = input.iconPath = '';
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.vesselCount = input.vesselCount = { min: 0 };
    this.name = input.name = '';
    this.onlyForVessels = input.onlyForVessels = false;
    this.slug = input.slug = '';
    this.temperatureRequired = input.temperatureRequired = false;
    this.yieldsNothing = input.yieldsNothing = false;
  }
}
