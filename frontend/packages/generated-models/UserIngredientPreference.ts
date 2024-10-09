// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IUserIngredientPreference {
  belongsToUser: string;
  createdAt: string;
  ingredient: ValidIngredient;
  notes: string;
  rating: number;
  archivedAt?: string;
  id: string;
  lastUpdatedAt?: string;
  allergy: boolean;
}

export class UserIngredientPreference implements IUserIngredientPreference {
  belongsToUser: string;
  createdAt: string;
  ingredient: ValidIngredient;
  notes: string;
  rating: number;
  archivedAt?: string;
  id: string;
  lastUpdatedAt?: string;
  allergy: boolean;
  constructor(input: Partial<UserIngredientPreference> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.ingredient = input.ingredient = new ValidIngredient();
    this.notes = input.notes = '';
    this.rating = input.rating = 0;
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.allergy = input.allergy = false;
  }
}
