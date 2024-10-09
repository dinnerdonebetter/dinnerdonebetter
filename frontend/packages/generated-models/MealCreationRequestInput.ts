// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponentCreationRequestInput } from './MealComponentCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealCreationRequestInput {
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
}

export class MealCreationRequestInput implements IMealCreationRequestInput {
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  constructor(input: Partial<MealCreationRequestInput> = {}) {
    this.name = input.name = '';
    this.components = input.components = new MealComponentCreationRequestInput();
    this.description = input.description = '';
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
  }
}
