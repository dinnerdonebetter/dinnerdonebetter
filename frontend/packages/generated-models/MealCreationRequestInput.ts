// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealComponentCreationRequestInput } from './MealComponentCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealCreationRequestInput {
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
  eligibleForMealPlans: boolean;
}

export class MealCreationRequestInput implements IMealCreationRequestInput {
  estimatedPortions: NumberRangeWithOptionalMax;
  name: string;
  components: MealComponentCreationRequestInput;
  description: string;
  eligibleForMealPlans: boolean;
  constructor(input: Partial<MealCreationRequestInput> = {}) {
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.name = input.name = '';
    this.components = input.components = new MealComponentCreationRequestInput();
    this.description = input.description = '';
    this.eligibleForMealPlans = input.eligibleForMealPlans = false;
  }
}
