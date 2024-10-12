// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredientStateAttributeType } from './enums';


export interface IValidIngredientState {
   archivedAt: string;
 attributeType: ValidIngredientStateAttributeType;
 createdAt: string;
 description: string;
 iconPath: string;
 id: string;
 lastUpdatedAt: string;
 name: string;
 pastTense: string;
 slug: string;

}

export class ValidIngredientState implements IValidIngredientState {
   archivedAt: string;
 attributeType: ValidIngredientStateAttributeType;
 createdAt: string;
 description: string;
 iconPath: string;
 id: string;
 lastUpdatedAt: string;
 name: string;
 pastTense: string;
 slug: string;
constructor(input: Partial<ValidIngredientState> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.attributeType = input.attributeType || 'other';
 this.createdAt = input.createdAt || '';
 this.description = input.description || '';
 this.iconPath = input.iconPath || '';
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.name = input.name || '';
 this.pastTense = input.pastTense || '';
 this.slug = input.slug || '';
}
}