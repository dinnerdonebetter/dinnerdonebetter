/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIError } from '../models/APIError';
import type { APIResponse } from '../models/APIResponse';
import type { Recipe } from '../models/Recipe';
import type { RecipeCreationRequestInput } from '../models/RecipeCreationRequestInput';
import type { RecipePrepTask } from '../models/RecipePrepTask';
import type { RecipePrepTaskCreationRequestInput } from '../models/RecipePrepTaskCreationRequestInput';
import type { RecipePrepTaskStep } from '../models/RecipePrepTaskStep';
import type { RecipePrepTaskUpdateRequestInput } from '../models/RecipePrepTaskUpdateRequestInput';
import type { RecipeRating } from '../models/RecipeRating';
import type { RecipeRatingCreationRequestInput } from '../models/RecipeRatingCreationRequestInput';
import type { RecipeRatingUpdateRequestInput } from '../models/RecipeRatingUpdateRequestInput';
import type { RecipeStep } from '../models/RecipeStep';
import type { RecipeStepCompletionCondition } from '../models/RecipeStepCompletionCondition';
import type { RecipeStepCompletionConditionForExistingRecipeCreationRequestInput } from '../models/RecipeStepCompletionConditionForExistingRecipeCreationRequestInput';
import type { RecipeStepCompletionConditionUpdateRequestInput } from '../models/RecipeStepCompletionConditionUpdateRequestInput';
import type { RecipeStepCreationRequestInput } from '../models/RecipeStepCreationRequestInput';
import type { RecipeStepIngredient } from '../models/RecipeStepIngredient';
import type { RecipeStepIngredientCreationRequestInput } from '../models/RecipeStepIngredientCreationRequestInput';
import type { RecipeStepIngredientUpdateRequestInput } from '../models/RecipeStepIngredientUpdateRequestInput';
import type { RecipeStepInstrument } from '../models/RecipeStepInstrument';
import type { RecipeStepInstrumentCreationRequestInput } from '../models/RecipeStepInstrumentCreationRequestInput';
import type { RecipeStepInstrumentUpdateRequestInput } from '../models/RecipeStepInstrumentUpdateRequestInput';
import type { RecipeStepProduct } from '../models/RecipeStepProduct';
import type { RecipeStepProductCreationRequestInput } from '../models/RecipeStepProductCreationRequestInput';
import type { RecipeStepProductUpdateRequestInput } from '../models/RecipeStepProductUpdateRequestInput';
import type { RecipeStepUpdateRequestInput } from '../models/RecipeStepUpdateRequestInput';
import type { RecipeStepVessel } from '../models/RecipeStepVessel';
import type { RecipeStepVesselCreationRequestInput } from '../models/RecipeStepVesselCreationRequestInput';
import type { RecipeStepVesselUpdateRequestInput } from '../models/RecipeStepVesselUpdateRequestInput';
import type { RecipeUpdateRequestInput } from '../models/RecipeUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipesService {
  /**
   * Operation for fetching Recipe
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @returns any
   * @throws ApiError
   */
  public static getRecipes(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
  ): CancelablePromise<
    APIResponse & {
      data?: Array<Recipe>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes',
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating Recipe
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipe(requestBody: RecipeCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: Recipe;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving Recipe
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipe(recipeId: string): CancelablePromise<
    APIResponse & {
      data?: Recipe;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for fetching Recipe
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipe(recipeId: string): CancelablePromise<
    APIResponse & {
      data?: Recipe;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for updating Recipe
   * @param recipeId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipe(
    recipeId: string,
    requestBody: RecipeUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: Recipe;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}',
      path: {
        recipeID: recipeId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating Recipe
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static cloneRecipe(recipeId: string): CancelablePromise<
    APIResponse & {
      data?: Recipe;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/clone',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for fetching APIError
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeDag(recipeId: string): CancelablePromise<
    APIResponse & {
      data?: APIError;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/dag',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for creating
   * @param recipeId
   * @throws ApiError
   */
  public static uploadMediaForRecipe(recipeId: string): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/images',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for fetching
   * @param recipeId
   * @throws ApiError
   */
  public static getMermaidDiagramForRecipe(recipeId: string): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/mermaid',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for fetching RecipePrepTaskStep
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeMealPlanTasks(recipeId: string): CancelablePromise<
    APIResponse & {
      data?: RecipePrepTaskStep;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/prep_steps',
      path: {
        recipeID: recipeId,
      },
    });
  }
  /**
   * Operation for fetching RecipePrepTask
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipePrepTasks(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipePrepTask>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/prep_tasks',
      path: {
        recipeID: recipeId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipePrepTask
   * @param recipeId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipePrepTask(
    recipeId: string,
    requestBody: RecipePrepTaskCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipePrepTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/prep_tasks',
      path: {
        recipeID: recipeId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipePrepTask
   * @param recipeId
   * @param recipePrepTaskId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipePrepTask(
    recipeId: string,
    recipePrepTaskId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipePrepTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}',
      path: {
        recipeID: recipeId,
        recipePrepTaskID: recipePrepTaskId,
      },
    });
  }
  /**
   * Operation for fetching RecipePrepTask
   * @param recipeId
   * @param recipePrepTaskId
   * @returns any
   * @throws ApiError
   */
  public static getRecipePrepTask(
    recipeId: string,
    recipePrepTaskId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipePrepTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}',
      path: {
        recipeID: recipeId,
        recipePrepTaskID: recipePrepTaskId,
      },
    });
  }
  /**
   * Operation for updating RecipePrepTask
   * @param recipeId
   * @param recipePrepTaskId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipePrepTask(
    recipeId: string,
    recipePrepTaskId: string,
    requestBody: RecipePrepTaskUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipePrepTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}',
      path: {
        recipeID: recipeId,
        recipePrepTaskID: recipePrepTaskId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeRating
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeRatings(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeRating>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/ratings',
      path: {
        recipeID: recipeId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeRating
   * @param recipeId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeRating(
    recipeId: string,
    requestBody: RecipeRatingCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/ratings',
      path: {
        recipeID: recipeId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeRating(
    recipeId: string,
    recipeRatingId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
    });
  }
  /**
   * Operation for fetching RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeRating(
    recipeId: string,
    recipeRatingId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
    });
  }
  /**
   * Operation for updating RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeRating(
    recipeId: string,
    recipeRatingId: string,
    requestBody: RecipeRatingUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStep
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeSteps(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStep>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps',
      path: {
        recipeID: recipeId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStep
   * @param recipeId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStep(
    recipeId: string,
    requestBody: RecipeStepCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStep;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps',
      path: {
        recipeID: recipeId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStep
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStep(
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStep;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStep
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStep(
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStep;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
    });
  }
  /**
   * Operation for updating RecipeStep
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStep(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStep;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStepCompletionCondition
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepCompletionConditions(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepCompletionCondition>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
    requestBody: RecipeStepCompletionConditionUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating
   * @param recipeId
   * @param recipeStepId
   * @throws ApiError
   */
  public static uploadMediaForRecipeStep(recipeId: string, recipeStepId: string): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/images',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepIngredient
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepIngredients(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepIngredientCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
    requestBody: RecipeStepIngredientUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStepInstrument
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepInstruments(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepInstrumentCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
    requestBody: RecipeStepInstrumentUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStepProduct
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepProducts(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepProduct>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepProductCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
    requestBody: RecipeStepProductUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStepVessel
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepVessels(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for creating RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepVesselCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
    requestBody: RecipeStepVesselUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching Recipe
   * @param q the search query parameter
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @returns any
   * @throws ApiError
   */
  public static searchForRecipes(
    q: string,
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
  ): CancelablePromise<
    APIResponse & {
      data?: Array<Recipe>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/search',
      query: {
        q: q,
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
}
