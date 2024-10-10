// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';

export interface IUserIngredientPreference {
  allergy: boolean;
  archivedAt: string;
  belongsToUser: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  notes: string;
  rating: number;
}

export class UserIngredientPreference implements IUserIngredientPreference {
  allergy: boolean;
  archivedAt: string;
  belongsToUser: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  notes: string;
  rating: number;
  constructor(input: Partial<UserIngredientPreference> = {}) {
    this.allergy = input.allergy || false;
    this.archivedAt = input.archivedAt || '';
    this.belongsToUser = input.belongsToUser || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.rating = input.rating || 0;
  }
}
