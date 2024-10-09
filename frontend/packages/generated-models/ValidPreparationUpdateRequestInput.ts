// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidPreparationUpdateRequestInput {
  instrumentCount: OptionalNumberRange;
  name?: string;
  onlyForVessels?: boolean;
  slug?: string;
  yieldsNothing?: boolean;
  iconPath?: string;
  timeEstimateRequired?: boolean;
  vesselCount: OptionalNumberRange;
  consumesVessel?: boolean;
  description?: string;
  pastTense?: string;
  temperatureRequired?: boolean;
  conditionExpressionRequired?: boolean;
  ingredientCount: OptionalNumberRange;
  restrictToIngredients?: boolean;
}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
  instrumentCount: OptionalNumberRange;
  name?: string;
  onlyForVessels?: boolean;
  slug?: string;
  yieldsNothing?: boolean;
  iconPath?: string;
  timeEstimateRequired?: boolean;
  vesselCount: OptionalNumberRange;
  consumesVessel?: boolean;
  description?: string;
  pastTense?: string;
  temperatureRequired?: boolean;
  conditionExpressionRequired?: boolean;
  ingredientCount: OptionalNumberRange;
  restrictToIngredients?: boolean;
  constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
    this.instrumentCount = input.instrumentCount = {};
    this.name = input.name;
    this.onlyForVessels = input.onlyForVessels;
    this.slug = input.slug;
    this.yieldsNothing = input.yieldsNothing;
    this.iconPath = input.iconPath;
    this.timeEstimateRequired = input.timeEstimateRequired;
    this.vesselCount = input.vesselCount = {};
    this.consumesVessel = input.consumesVessel;
    this.description = input.description;
    this.pastTense = input.pastTense;
    this.temperatureRequired = input.temperatureRequired;
    this.conditionExpressionRequired = input.conditionExpressionRequired;
    this.ingredientCount = input.ingredientCount = {};
    this.restrictToIngredients = input.restrictToIngredients;
  }
}
