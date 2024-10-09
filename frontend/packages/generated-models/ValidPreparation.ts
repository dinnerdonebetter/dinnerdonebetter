// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IValidPreparation {
   description: string;
 instrumentCount: NumberRangeWithOptionalMax;
 name: string;
 pastTense: string;
 restrictToIngredients: boolean;
 slug: string;
 vesselCount: NumberRangeWithOptionalMax;
 createdAt: string;
 ingredientCount: NumberRangeWithOptionalMax;
 lastUpdatedAt?: string;
 timeEstimateRequired: boolean;
 iconPath: string;
 conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 id: string;
 onlyForVessels: boolean;
 archivedAt?: string;
 yieldsNothing: boolean;
 temperatureRequired: boolean;

}

export class ValidPreparation implements IValidPreparation {
   description: string;
 instrumentCount: NumberRangeWithOptionalMax;
 name: string;
 pastTense: string;
 restrictToIngredients: boolean;
 slug: string;
 vesselCount: NumberRangeWithOptionalMax;
 createdAt: string;
 ingredientCount: NumberRangeWithOptionalMax;
 lastUpdatedAt?: string;
 timeEstimateRequired: boolean;
 iconPath: string;
 conditionExpressionRequired: boolean;
 consumesVessel: boolean;
 id: string;
 onlyForVessels: boolean;
 archivedAt?: string;
 yieldsNothing: boolean;
 temperatureRequired: boolean;
constructor(input: Partial<ValidPreparation> = {}) {
	 this.description = input.description = '';
 this.instrumentCount = input.instrumentCount = { min: 0 };
 this.name = input.name = '';
 this.pastTense = input.pastTense = '';
 this.restrictToIngredients = input.restrictToIngredients = false;
 this.slug = input.slug = '';
 this.vesselCount = input.vesselCount = { min: 0 };
 this.createdAt = input.createdAt = '';
 this.ingredientCount = input.ingredientCount = { min: 0 };
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.timeEstimateRequired = input.timeEstimateRequired = false;
 this.iconPath = input.iconPath = '';
 this.conditionExpressionRequired = input.conditionExpressionRequired = false;
 this.consumesVessel = input.consumesVessel = false;
 this.id = input.id = '';
 this.onlyForVessels = input.onlyForVessels = false;
 this.archivedAt = input.archivedAt;
 this.yieldsNothing = input.yieldsNothing = false;
 this.temperatureRequired = input.temperatureRequired = false;
}
}