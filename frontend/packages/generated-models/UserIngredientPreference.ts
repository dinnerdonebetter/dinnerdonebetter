// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IUserIngredientPreference {
  belongsToUser: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  rating: number;
  allergy: boolean;
  archivedAt?: string;
  ingredient: ValidIngredient;
  notes: string;
}

export class UserIngredientPreference implements IUserIngredientPreference {
  belongsToUser: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  rating: number;
  allergy: boolean;
  archivedAt?: string;
  ingredient: ValidIngredient;
  notes: string;
  constructor(input: Partial<UserIngredientPreference> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.rating = input.rating = 0;
    this.allergy = input.allergy = false;
    this.archivedAt = input.archivedAt;
    this.ingredient = input.ingredient = new ValidIngredient();
    this.notes = input.notes = '';
  }
}
