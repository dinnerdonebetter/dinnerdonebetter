// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IValidPreparationCreationRequestInput {
   name: string;
 temperatureRequired: boolean;
 vesselCount: NumberRangeWithOptionalMax;
 instrumentCount: NumberRangeWithOptionalMax;
 restrictToIngredients: boolean;
 timeEstimateRequired: boolean;
 yieldsNothing: boolean;
 ingredientCount: NumberRangeWithOptionalMax;
 slug: string;
 conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 description: string;
 iconPath: string;
 onlyForVessels: boolean;
 pastTense: string;

}

export class ValidPreparationCreationRequestInput implements IValidPreparationCreationRequestInput {
   name: string;
 temperatureRequired: boolean;
 vesselCount: NumberRangeWithOptionalMax;
 instrumentCount: NumberRangeWithOptionalMax;
 restrictToIngredients: boolean;
 timeEstimateRequired: boolean;
 yieldsNothing: boolean;
 ingredientCount: NumberRangeWithOptionalMax;
 slug: string;
 conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 description: string;
 iconPath: string;
 onlyForVessels: boolean;
 pastTense: string;
constructor(input: Partial<ValidPreparationCreationRequestInput> = {}) {
	 this.name = input.name = '';
 this.temperatureRequired = input.temperatureRequired = false;
 this.vesselCount = input.vesselCount = { min: 0 };
 this.instrumentCount = input.instrumentCount = { min: 0 };
 this.restrictToIngredients = input.restrictToIngredients = false;
 this.timeEstimateRequired = input.timeEstimateRequired = false;
 this.yieldsNothing = input.yieldsNothing = false;
 this.ingredientCount = input.ingredientCount = { min: 0 };
 this.slug = input.slug = '';
 this.conditionExpressionRequired = input.conditionExpressionRequired = false;
 this.consumesVessel = input.consumesVessel = false;
 this.description = input.description = '';
 this.iconPath = input.iconPath = '';
 this.onlyForVessels = input.onlyForVessels = false;
 this.pastTense = input.pastTense = '';
}
}