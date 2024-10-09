// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidPreparationUpdateRequestInput {
  description?: string;
  onlyForVessels?: boolean;
  restrictToIngredients?: boolean;
  instrumentCount: OptionalNumberRange;
  timeEstimateRequired?: boolean;
  iconPath?: string;
  name?: string;
  pastTense?: string;
  temperatureRequired?: boolean;
  vesselCount: OptionalNumberRange;
  conditionExpressionRequired?: boolean;
  consumesVessel?: boolean;
  ingredientCount: OptionalNumberRange;
  slug?: string;
  yieldsNothing?: boolean;
}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
  description?: string;
  onlyForVessels?: boolean;
  restrictToIngredients?: boolean;
  instrumentCount: OptionalNumberRange;
  timeEstimateRequired?: boolean;
  iconPath?: string;
  name?: string;
  pastTense?: string;
  temperatureRequired?: boolean;
  vesselCount: OptionalNumberRange;
  conditionExpressionRequired?: boolean;
  consumesVessel?: boolean;
  ingredientCount: OptionalNumberRange;
  slug?: string;
  yieldsNothing?: boolean;
  constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
    this.description = input.description;
    this.onlyForVessels = input.onlyForVessels;
    this.restrictToIngredients = input.restrictToIngredients;
    this.instrumentCount = input.instrumentCount = {};
    this.timeEstimateRequired = input.timeEstimateRequired;
    this.iconPath = input.iconPath;
    this.name = input.name;
    this.pastTense = input.pastTense;
    this.temperatureRequired = input.temperatureRequired;
    this.vesselCount = input.vesselCount = {};
    this.conditionExpressionRequired = input.conditionExpressionRequired;
    this.consumesVessel = input.consumesVessel;
    this.ingredientCount = input.ingredientCount = {};
    this.slug = input.slug;
    this.yieldsNothing = input.yieldsNothing;
  }
}
