// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidPreparationUpdateRequestInput {
  yieldsNothing?: boolean;
  conditionExpressionRequired?: boolean;
  ingredientCount: OptionalNumberRange;
  onlyForVessels?: boolean;
  pastTense?: string;
  timeEstimateRequired?: boolean;
  restrictToIngredients?: boolean;
  consumesVessel?: boolean;
  description?: string;
  iconPath?: string;
  instrumentCount: OptionalNumberRange;
  slug?: string;
  name?: string;
  temperatureRequired?: boolean;
  vesselCount: OptionalNumberRange;
}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
  yieldsNothing?: boolean;
  conditionExpressionRequired?: boolean;
  ingredientCount: OptionalNumberRange;
  onlyForVessels?: boolean;
  pastTense?: string;
  timeEstimateRequired?: boolean;
  restrictToIngredients?: boolean;
  consumesVessel?: boolean;
  description?: string;
  iconPath?: string;
  instrumentCount: OptionalNumberRange;
  slug?: string;
  name?: string;
  temperatureRequired?: boolean;
  vesselCount: OptionalNumberRange;
  constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
    this.yieldsNothing = input.yieldsNothing;
    this.conditionExpressionRequired = input.conditionExpressionRequired;
    this.ingredientCount = input.ingredientCount = {};
    this.onlyForVessels = input.onlyForVessels;
    this.pastTense = input.pastTense;
    this.timeEstimateRequired = input.timeEstimateRequired;
    this.restrictToIngredients = input.restrictToIngredients;
    this.consumesVessel = input.consumesVessel;
    this.description = input.description;
    this.iconPath = input.iconPath;
    this.instrumentCount = input.instrumentCount = {};
    this.slug = input.slug;
    this.name = input.name;
    this.temperatureRequired = input.temperatureRequired;
    this.vesselCount = input.vesselCount = {};
  }
}
