// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IValidPreparationUpdateRequestInput {
   conditionExpressionRequired?: boolean;
 pastTense?: string;
 slug?: string;
 temperatureRequired?: boolean;
 description?: string;
 instrumentCount: OptionalNumberRange;
 iconPath?: string;
 ingredientCount: OptionalNumberRange;
 name?: string;
 restrictToIngredients?: boolean;
 yieldsNothing?: boolean;
 consumesVessel?: boolean;
 onlyForVessels?: boolean;
 timeEstimateRequired?: boolean;
 vesselCount: OptionalNumberRange;

}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
   conditionExpressionRequired?: boolean;
 pastTense?: string;
 slug?: string;
 temperatureRequired?: boolean;
 description?: string;
 instrumentCount: OptionalNumberRange;
 iconPath?: string;
 ingredientCount: OptionalNumberRange;
 name?: string;
 restrictToIngredients?: boolean;
 yieldsNothing?: boolean;
 consumesVessel?: boolean;
 onlyForVessels?: boolean;
 timeEstimateRequired?: boolean;
 vesselCount: OptionalNumberRange;
constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
	 this.conditionExpressionRequired = input.conditionExpressionRequired;
 this.pastTense = input.pastTense;
 this.slug = input.slug;
 this.temperatureRequired = input.temperatureRequired;
 this.description = input.description;
 this.instrumentCount = input.instrumentCount = {};
 this.iconPath = input.iconPath;
 this.ingredientCount = input.ingredientCount = {};
 this.name = input.name;
 this.restrictToIngredients = input.restrictToIngredients;
 this.yieldsNothing = input.yieldsNothing;
 this.consumesVessel = input.consumesVessel;
 this.onlyForVessels = input.onlyForVessels;
 this.timeEstimateRequired = input.timeEstimateRequired;
 this.vesselCount = input.vesselCount = {};
}
}