// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IValidPreparationUpdateRequestInput {
   conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 description: string;
 iconPath: string;
 ingredientCount: OptionalNumberRange;
 instrumentCount: OptionalNumberRange;
 name: string;
 onlyForVessels: boolean;
 pastTense: string;
 restrictToIngredients: boolean;
 slug: string;
 temperatureRequired: boolean;
 timeEstimateRequired: boolean;
 vesselCount: OptionalNumberRange;
 yieldsNothing: boolean;

}

export class ValidPreparationUpdateRequestInput implements IValidPreparationUpdateRequestInput {
   conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 description: string;
 iconPath: string;
 ingredientCount: OptionalNumberRange;
 instrumentCount: OptionalNumberRange;
 name: string;
 onlyForVessels: boolean;
 pastTense: string;
 restrictToIngredients: boolean;
 slug: string;
 temperatureRequired: boolean;
 timeEstimateRequired: boolean;
 vesselCount: OptionalNumberRange;
 yieldsNothing: boolean;
constructor(input: Partial<ValidPreparationUpdateRequestInput> = {}) {
	 this.conditionExpressionRequired = input.conditionExpressionRequired || false;
 this.consumesVessel = input.consumesVessel || false;
 this.description = input.description || '';
 this.iconPath = input.iconPath || '';
 this.ingredientCount = input.ingredientCount || {};
 this.instrumentCount = input.instrumentCount || {};
 this.name = input.name || '';
 this.onlyForVessels = input.onlyForVessels || false;
 this.pastTense = input.pastTense || '';
 this.restrictToIngredients = input.restrictToIngredients || false;
 this.slug = input.slug || '';
 this.temperatureRequired = input.temperatureRequired || false;
 this.timeEstimateRequired = input.timeEstimateRequired || false;
 this.vesselCount = input.vesselCount || {};
 this.yieldsNothing = input.yieldsNothing || false;
}
}