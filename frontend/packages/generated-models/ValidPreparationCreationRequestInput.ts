// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidPreparationCreationRequestInput {
  description: string;
  instrumentCount: NumberRangeWithOptionalMax;
  pastTense: string;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  timeEstimateRequired: boolean;
  iconPath: string;
  name: string;
  restrictToIngredients: boolean;
  yieldsNothing: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  ingredientCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  slug: string;
  temperatureRequired: boolean;
}

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
  description: string;
  instrumentCount: NumberRangeWithOptionalMax;
  pastTense: string;
  conditionExpressionRequired: boolean;
  consumesVessel: boolean;
  timeEstimateRequired: boolean;
  iconPath: string;
  name: string;
  restrictToIngredients: boolean;
  yieldsNothing: boolean;
  vesselCount: NumberRangeWithOptionalMax;
  ingredientCount: NumberRangeWithOptionalMax;
  onlyForVessels: boolean;
  slug: string;
  temperatureRequired: boolean;
  constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
    this.description = input.description = '';
    this.instrumentCount = input.instrumentCount = { min: 0 };
    this.pastTense = input.pastTense = '';
    this.conditionExpressionRequired = input.conditionExpressionRequired = false;
    this.consumesVessel = input.consumesVessel = false;
    this.timeEstimateRequired = input.timeEstimateRequired = false;
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.restrictToIngredients = input.restrictToIngredients = false;
    this.yieldsNothing = input.yieldsNothing = false;
    this.vesselCount = input.vesselCount = { min: 0 };
    this.ingredientCount = input.ingredientCount = { min: 0 };
    this.onlyForVessels = input.onlyForVessels = false;
    this.slug = input.slug = '';
    this.temperatureRequired = input.temperatureRequired = false;
  }
}
