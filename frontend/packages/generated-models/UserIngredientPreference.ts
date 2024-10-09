// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredient } from './ValidIngredient';


export interface IUserIngredientPreference {
   notes: string;
 rating: number;
 archivedAt?: string;
 ingredient: ValidIngredient;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 allergy: boolean;
 belongsToUser: string;

}

export class UserIngredientPreference implements IUserIngredientPreference {
   notes: string;
 rating: number;
 archivedAt?: string;
 ingredient: ValidIngredient;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 allergy: boolean;
 belongsToUser: string;
constructor(input: Partial<UserIngredientPreference> = {}) {
	 this.notes = input.notes = '';
 this.rating = input.rating = 0;
 this.archivedAt = input.archivedAt;
 this.ingredient = input.ingredient = new ValidIngredient();
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.allergy = input.allergy = false;
 this.belongsToUser = input.belongsToUser = '';
}
}