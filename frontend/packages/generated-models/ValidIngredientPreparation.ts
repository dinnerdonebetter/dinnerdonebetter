// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredient } from './ValidIngredient';
 import { ValidPreparation } from './ValidPreparation';


export interface IValidIngredientPreparation {
   notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;
 ingredient: ValidIngredient;
 lastUpdatedAt?: string;

}

export class ValidIngredientPreparation implements IValidIngredientPreparation {
   notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;
 ingredient: ValidIngredient;
 lastUpdatedAt?: string;
constructor(input: Partial<ValidIngredientPreparation> = {}) {
	 this.notes = input.notes = '';
 this.preparation = input.preparation = new ValidPreparation();
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.ingredient = input.ingredient = new ValidIngredient();
 this.lastUpdatedAt = input.lastUpdatedAt;
}
}