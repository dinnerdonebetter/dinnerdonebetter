// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IUserIngredientPreference {
  belongsToUser: string;
  id: string;
  ingredient: ValidIngredient;
  allergy: boolean;
  archivedAt?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  notes: string;
  rating: number;
}

export class UserIngredientPreference implements IUserIngredientPreference {
  belongsToUser: string;
  id: string;
  ingredient: ValidIngredient;
  allergy: boolean;
  archivedAt?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  notes: string;
  rating: number;
  constructor(input: Partial<UserIngredientPreference> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.id = input.id = '';
    this.ingredient = input.ingredient = new ValidIngredient();
    this.allergy = input.allergy = false;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.rating = input.rating = 0;
  }
}
