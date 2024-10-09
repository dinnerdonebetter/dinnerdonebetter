// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredientStateAttributeType } from './enums';


export interface IValidIngredientState {
   attributeType: ValidIngredientStateAttributeType;
 createdAt: string;
 iconPath: string;
 name: string;
 slug: string;
 archivedAt?: string;
 description: string;
 id: string;
 lastUpdatedAt?: string;
 pastTense: string;

}

export class ValidIngredientState implements IValidIngredientState {
   attributeType: ValidIngredientStateAttributeType;
 createdAt: string;
 iconPath: string;
 name: string;
 slug: string;
 archivedAt?: string;
 description: string;
 id: string;
 lastUpdatedAt?: string;
 pastTense: string;
constructor(input: Partial<ValidIngredientState> = {}) {
	 this.attributeType = input.attributeType = 'other';
 this.createdAt = input.createdAt = '';
 this.iconPath = input.iconPath = '';
 this.name = input.name = '';
 this.slug = input.slug = '';
 this.archivedAt = input.archivedAt;
 this.description = input.description = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.pastTense = input.pastTense = '';
}
}