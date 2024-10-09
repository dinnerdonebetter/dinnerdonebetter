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
  temperatureInCelsius: NumberRange;
  archivedAt?: string;
  conditionExpression: string;
  instruments: RecipeStepInstrument;
  startTimerAutomatically: boolean;
  completionConditions: RecipeStepCompletionCondition;
  explicitInstructions: string;
  lastUpdatedAt?: string;
  optional: boolean;
  media: RecipeMedia;
  notes: string;
  vessels: RecipeStepVessel;
  belongsToRecipe: string;
  estimatedTimeInSeconds: NumberRange;
  id: string;
  ingredients: RecipeStepIngredient;
  createdAt: string;
  index: number;
  preparation: ValidPreparation;
  products: RecipeStepProduct;
}

export class RecipeStep implements IRecipeStep {
  temperatureInCelsius: NumberRange;
  archivedAt?: string;
  conditionExpression: string;
  instruments: RecipeStepInstrument;
  startTimerAutomatically: boolean;
  completionConditions: RecipeStepCompletionCondition;
  explicitInstructions: string;
  lastUpdatedAt?: string;
  optional: boolean;
  media: RecipeMedia;
  notes: string;
  vessels: RecipeStepVessel;
  belongsToRecipe: string;
  estimatedTimeInSeconds: NumberRange;
  id: string;
  ingredients: RecipeStepIngredient;
  createdAt: string;
  index: number;
  preparation: ValidPreparation;
  products: RecipeStepProduct;
  constructor(input: Partial<RecipeStep> = {}) {
    this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
    this.archivedAt = input.archivedAt;
    this.conditionExpression = input.conditionExpression = '';
    this.instruments = input.instruments = new RecipeStepInstrument();
    this.startTimerAutomatically = input.startTimerAutomatically = false;
    this.completionConditions = input.completionConditions = new RecipeStepCompletionCondition();
    this.explicitInstructions = input.explicitInstructions = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.optional = input.optional = false;
    this.media = input.media = new RecipeMedia();
    this.notes = input.notes = '';
    this.vessels = input.vessels = new RecipeStepVessel();
    this.belongsToRecipe = input.belongsToRecipe = '';
    this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
    this.id = input.id = '';
    this.ingredients = input.ingredients = new RecipeStepIngredient();
    this.createdAt = input.createdAt = '';
    this.index = input.index = 0;
    this.preparation = input.preparation = new ValidPreparation();
    this.products = input.products = new RecipeStepProduct();
  }
}
