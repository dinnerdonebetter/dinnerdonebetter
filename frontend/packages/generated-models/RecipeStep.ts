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
  media: RecipeMedia;
  notes: string;
  products: RecipeStepProduct;
  archivedAt?: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  id: string;
  ingredients: RecipeStepIngredient;
  belongsToRecipe: string;
  completionConditions: RecipeStepCompletionCondition;
  explicitInstructions: string;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  lastUpdatedAt?: string;
  optional: boolean;
  vessels: RecipeStepVessel;
  index: number;
  instruments: RecipeStepInstrument;
  temperatureInCelsius: NumberRange;
}

export class RecipeStep implements IRecipeStep {
  media: RecipeMedia;
  notes: string;
  products: RecipeStepProduct;
  archivedAt?: string;
  createdAt: string;
  estimatedTimeInSeconds: NumberRange;
  id: string;
  ingredients: RecipeStepIngredient;
  belongsToRecipe: string;
  completionConditions: RecipeStepCompletionCondition;
  explicitInstructions: string;
  preparation: ValidPreparation;
  startTimerAutomatically: boolean;
  conditionExpression: string;
  lastUpdatedAt?: string;
  optional: boolean;
  vessels: RecipeStepVessel;
  index: number;
  instruments: RecipeStepInstrument;
  temperatureInCelsius: NumberRange;
  constructor(input: Partial<RecipeStep> = {}) {
    this.media = input.media = new RecipeMedia();
    this.notes = input.notes = '';
    this.products = input.products = new RecipeStepProduct();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.id = input.id = '';
    this.ingredients = input.ingredients = new RecipeStepIngredient();
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.completionConditions = input.completionConditions = new RecipeStepCompletionCondition();
    this.explicitInstructions = input.explicitInstructions = '';
    this.preparation = input.preparation = new ValidPreparation();
    this.startTimerAutomatically = input.startTimerAutomatically = false;
    this.conditionExpression = input.conditionExpression = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.optional = input.optional = false;
    this.vessels = input.vessels = new RecipeStepVessel();
    this.index = input.index = 0;
    this.instruments = input.instruments = new RecipeStepInstrument();
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
  }
}
