// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredientStateAttributeType } from './enums';


export interface IValidIngredientStateUpdateRequestInput {
   name?: string;
 pastTense?: string;
 slug?: string;
 attributeType?: ValidIngredientStateAttributeType;
 description?: string;
 iconPath?: string;

}

export class ValidIngredientStateUpdateRequestInput implements IValidIngredientStateUpdateRequestInput {
   name?: string;
 pastTense?: string;
 slug?: string;
 attributeType?: ValidIngredientStateAttributeType;
 description?: string;
 iconPath?: string;
constructor(input: Partial<ValidIngredientStateUpdateRequestInput> = {}) {
	 this.name = input.name;
 this.pastTense = input.pastTense;
 this.slug = input.slug;
 this.attributeType = input.attributeType;
 this.description = input.description;
 this.iconPath = input.iconPath;
}
}