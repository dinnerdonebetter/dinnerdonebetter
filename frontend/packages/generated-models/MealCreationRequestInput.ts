// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponentCreationRequestInput } from './MealComponentCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealCreationRequestInput {
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
}

export class MealCreationRequestInput implements IMealCreationRequestInput {
  eligibleForMealPlans: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
  constructor(input: Partial<MealCreationRequestInput> = {}) {
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.name = input.name = '';
    this.components = input.components = new MealComponentCreationRequestInput();
    this.description = input.description = '';
  }
}
