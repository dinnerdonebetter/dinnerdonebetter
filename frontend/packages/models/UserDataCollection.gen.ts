// GENERATED CODE, DO NOT EDIT MANUALLY

import { AuditLogEntry } from './AuditLogEntry.gen';
import { Household } from './Household.gen';
import { HouseholdInvitation } from './HouseholdInvitation.gen';
import { Meal } from './Meal.gen';
import { Recipe } from './Recipe.gen';
import { RecipeRating } from './RecipeRating.gen';
import { ServiceSettingConfiguration } from './ServiceSettingConfiguration.gen';
import { User } from './User.gen';
import { UserIngredientPreference } from './UserIngredientPreference.gen';

export interface IUserDataCollection {
  auditLogEntries: object;
  householdInstrumentOwnerships: object;
  households: Household[];
  mealPlans: object;
  meals: Meal[];
  receivedInvites: HouseholdInvitation[];
  recipeRatings: RecipeRating[];
  recipes: Recipe[];
  reportID: string;
  sentInvites: HouseholdInvitation[];
  serviceSettingConfigurations: object;
  user: User;
  userAuditLogEntries: AuditLogEntry[];
  userIngredientPreferences: UserIngredientPreference[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
}

export class UserDataCollection implements IUserDataCollection {
  auditLogEntries: object;
  householdInstrumentOwnerships: object;
  households: Household[];
  mealPlans: object;
  meals: Meal[];
  receivedInvites: HouseholdInvitation[];
  recipeRatings: RecipeRating[];
  recipes: Recipe[];
  reportID: string;
  sentInvites: HouseholdInvitation[];
  serviceSettingConfigurations: object;
  user: User;
  userAuditLogEntries: AuditLogEntry[];
  userIngredientPreferences: UserIngredientPreference[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
  constructor(input: Partial<UserDataCollection> = {}) {
    this.auditLogEntries = input.auditLogEntries || {};
    this.householdInstrumentOwnerships = input.householdInstrumentOwnerships || {};
    this.households = input.households || [];
    this.mealPlans = input.mealPlans || {};
    this.meals = input.meals || [];
    this.receivedInvites = input.receivedInvites || [];
    this.recipeRatings = input.recipeRatings || [];
    this.recipes = input.recipes || [];
    this.reportID = input.reportID || '';
    this.sentInvites = input.sentInvites || [];
    this.serviceSettingConfigurations = input.serviceSettingConfigurations || {};
    this.user = input.user || new User();
    this.userAuditLogEntries = input.userAuditLogEntries || [];
    this.userIngredientPreferences = input.userIngredientPreferences || [];
    this.userServiceSettingConfigurations = input.userServiceSettingConfigurations || [];
    this.webhooks = input.webhooks || {};
  }
}
