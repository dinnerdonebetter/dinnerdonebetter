// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia.gen';
import { RecipeStepCompletionCondition } from './RecipeStepCompletionCondition.gen';
import { RecipeStepIngredient } from './RecipeStepIngredient.gen';
import { RecipeStepInstrument } from './RecipeStepInstrument.gen';
import { RecipeStepProduct } from './RecipeStepProduct.gen';
import { RecipeStepVessel } from './RecipeStepVessel.gen';
import { ValidPreparation } from './ValidPreparation.gen';
import { NumberRange } from './number_range.gen';

export interface IRecipeStep {
  archivedAt: string;
  belongsToRecipe: string;
  completionConditions: RecipeStepCompletionCondition[];
  conditionExpression: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  id: string;
  index: number;
  ingredients: RecipeStepIngredient[];
  instruments: RecipeStepInstrument[];
  lastUpdatedAt: string;
  media: RecipeMedia[];
  notes: string;
  optional: boolean;
  preparation: ValidPreparation;
  products: RecipeStepProduct[];
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVessel[];
}

export class RecipeStep implements IRecipeStep {
  archivedAt: string;
  belongsToRecipe: string;
  completionConditions: RecipeStepCompletionCondition[];
  conditionExpression: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  explicitInstructions: string;
  id: string;
  index: number;
  ingredients: RecipeStepIngredient[];
  instruments: RecipeStepInstrument[];
  lastUpdatedAt: string;
  media: RecipeMedia[];
  notes: string;
  optional: boolean;
  preparation: ValidPreparation;
  products: RecipeStepProduct[];
  startTimerAutomatically: boolean;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVessel[];
  constructor(input: Partial<RecipeStep> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.completionConditions = input.completionConditions || [];
    this.conditionExpression = input.conditionExpression || '';
    this.createdAt = input.createdAt || '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds || { min: 0, max: 0 };
    this.explicitInstructions = input.explicitInstructions || '';
    this.id = input.id || '';
    this.index = input.index || 0;
    this.ingredients = input.ingredients || [];
    this.instruments = input.instruments || [];
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.media = input.media || [];
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.preparation = input.preparation || new ValidPreparation();
    this.products = input.products || [];
    this.startTimerAutomatically = input.startTimerAutomatically || false;
    this.temperatureInCelsius = input.temperatureInCelsius || { min: 0, max: 0 };
    this.vessels = input.vessels || [];
  }
}
