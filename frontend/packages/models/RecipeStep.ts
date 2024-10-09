// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipeStepCompletionCondition } from './RecipeStepCompletionCondition';
import { RecipeStepIngredient } from './RecipeStepIngredient';
import { RecipeStepInstrument } from './RecipeStepInstrument';
import { RecipeStepProduct } from './RecipeStepProduct';
import { RecipeStepVessel } from './RecipeStepVessel';
import { ValidPreparation } from './ValidPreparation';
import { NumberRange } from './number_range';

export interface IRecipeStep {
  notes: string;
  startTimerAutomatically: boolean;
  archivedAt: string;
  belongsToRecipe: string;
  explicitInstructions: string;
  id: string;
  instruments: RecipeStepInstrument[];
  media: RecipeMedia[];
  vessels: RecipeStepVessel[];
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  products: RecipeStepProduct[];
  temperatureInCelsius: NumberRange;
  index: number;
  lastUpdatedAt: string;
  completionConditions: RecipeStepCompletionCondition[];
  conditionExpression: string;
  ingredients: RecipeStepIngredient[];
  optional: boolean;
  preparation: ValidPreparation;
}

export class RecipeStep implements IRecipeStep {
  notes: string;
  startTimerAutomatically: boolean;
  archivedAt: string;
  belongsToRecipe: string;
  explicitInstructions: string;
  id: string;
  instruments: RecipeStepInstrument[];
  media: RecipeMedia[];
  vessels: RecipeStepVessel[];
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  products: RecipeStepProduct[];
  temperatureInCelsius: NumberRange;
  index: number;
  lastUpdatedAt: string;
  completionConditions: RecipeStepCompletionCondition[];
  conditionExpression: string;
  ingredients: RecipeStepIngredient[];
  optional: boolean;
  preparation: ValidPreparation;
  constructor(input: Partial<RecipeStep> = {}) {
    this.notes = input.notes || '';
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.explicitInstructions = input.explicitInstructions || '';
    this.id = input.id || '';
    this.instruments = input.instruments || [];
    this.media = input.media || [];
    this.vessels = input.vessels || [];
    this.createdAt = input.createdAt || '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.products = input.products || [];
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
    this.index = input.index || 0;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.completionConditions = input.completionConditions || [];
    this.conditionExpression = input.conditionExpression || '';
    this.ingredients = input.ingredients || [];
    this.optional = input.optional || false;
    this.preparation = input.preparation || new ValidPreparation();
  }
}
