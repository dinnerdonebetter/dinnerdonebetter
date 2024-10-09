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
  preparation: ValidPreparation;
  completionConditions: RecipeStepCompletionCondition;
  conditionExpression: string;
  explicitInstructions: string;
  index: number;
  optional: boolean;
  startTimerAutomatically: boolean;
  lastUpdatedAt?: string;
  media: RecipeMedia;
  products: RecipeStepProduct;
  belongsToRecipe: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  instruments: RecipeStepInstrument;
  notes: string;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVessel;
  archivedAt?: string;
  id: string;
  ingredients: RecipeStepIngredient;
}

export class RecipeStep implements IRecipeStep {
  preparation: ValidPreparation;
  completionConditions: RecipeStepCompletionCondition;
  conditionExpression: string;
  explicitInstructions: string;
  index: number;
  optional: boolean;
  startTimerAutomatically: boolean;
  lastUpdatedAt?: string;
  media: RecipeMedia;
  products: RecipeStepProduct;
  belongsToRecipe: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  instruments: RecipeStepInstrument;
  notes: string;
  temperatureInCelsius: NumberRange;
  vessels: RecipeStepVessel;
  archivedAt?: string;
  id: string;
  ingredients: RecipeStepIngredient;
  constructor(input: Partial<RecipeStep> = {}) {
    this.preparation = input.preparation = new ValidPreparation();
    this.completionConditions = input.completionConditions = new RecipeStepCompletionCondition();
    this.conditionExpression = input.conditionExpression = '';
    this.explicitInstructions = input.explicitInstructions = '';
    this.index = input.index = 0;
    this.optional = input.optional = false;
    this.startTimerAutomatically = input.startTimerAutomatically = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.media = input.media = new RecipeMedia();
    this.products = input.products = new RecipeStepProduct();
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.createdAt = input.createdAt = '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.instruments = input.instruments = new RecipeStepInstrument();
    this.notes = input.notes = '';
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.vessels = input.vessels = new RecipeStepVessel();
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.ingredients = input.ingredients = new RecipeStepIngredient();
  }
}
