// Code generated by gen_typescript. DO NOT EDIT.

export interface IRecipeMedia {
  createdAt: NonNullable<string>;
  archivedAt?: string;
  lastUpdatedAt?: string;
  id: NonNullable<string>;
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType: NonNullable<string>;
  internalPath: NonNullable<string>;
  externalPath: NonNullable<string>;
  index: NonNullable<number>;
}

export class RecipeMedia implements IRecipeMedia {
  createdAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  archivedAt?: string;
  lastUpdatedAt?: string;
  id: NonNullable<string> = '';
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType: NonNullable<string> = '';
  internalPath: NonNullable<string> = '';
  externalPath: NonNullable<string> = '';
  index: NonNullable<number> = 0;

  constructor(input: Partial<RecipeMedia> = {}) {
    this.createdAt = input.createdAt ?? '1970-01-01T00:00:00Z';
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.id = input.id ?? '';
    this.belongsToRecipe = input.belongsToRecipe;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.mimeType = input.mimeType ?? '';
    this.internalPath = input.internalPath ?? '';
    this.externalPath = input.externalPath ?? '';
    this.index = input.index ?? 0;
  }
}

export interface IRecipeMediaCreationRequestInput {
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType: NonNullable<string>;
  internalPath: NonNullable<string>;
  externalPath: NonNullable<string>;
  index: NonNullable<number>;
}

export class RecipeMediaCreationRequestInput implements IRecipeMediaCreationRequestInput {
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType: NonNullable<string> = '';
  internalPath: NonNullable<string> = '';
  externalPath: NonNullable<string> = '';
  index: NonNullable<number> = 0;

  constructor(input: Partial<RecipeMediaCreationRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.mimeType = input.mimeType ?? '';
    this.internalPath = input.internalPath ?? '';
    this.externalPath = input.externalPath ?? '';
    this.index = input.index ?? 0;
  }
}

export interface IRecipeMediaUpdateRequestInput {
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType?: string;
  internalPath?: string;
  externalPath?: string;
  index?: number;
}

export class RecipeMediaUpdateRequestInput implements IRecipeMediaUpdateRequestInput {
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  mimeType?: string;
  internalPath?: string;
  externalPath?: string;
  index?: number;

  constructor(input: Partial<RecipeMediaUpdateRequestInput> = {}) {
    this.belongsToRecipe = input.belongsToRecipe;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.mimeType = input.mimeType;
    this.internalPath = input.internalPath;
    this.externalPath = input.externalPath;
    this.index = input.index;
  }
}