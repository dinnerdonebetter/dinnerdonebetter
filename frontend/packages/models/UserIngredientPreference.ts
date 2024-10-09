// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IUserIngredientPreference {
  allergy: boolean;
  belongsToUser: string;
  id: string;
  lastUpdatedAt: string;
  archivedAt: string;
  createdAt: string;
  ingredient: ValidIngredient;
  notes: string;
  rating: number;
}

export class UserIngredientPreference implements IUserIngredientPreference {
  allergy: boolean;
  belongsToUser: string;
  id: string;
  lastUpdatedAt: string;
  archivedAt: string;
  createdAt: string;
  ingredient: ValidIngredient;
  notes: string;
  rating: number;
  constructor(input: Partial<UserIngredientPreference> = {}) {
    this.allergy = input.allergy || false;
    this.belongsToUser = input.belongsToUser || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.notes = input.notes || '';
    this.rating = input.rating || 0;
  }
}
