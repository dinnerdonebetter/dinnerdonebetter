/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipePrepTask } from '../models/RecipePrepTask';
import type { RecipePrepTaskCreationRequestInput } from '../models/RecipePrepTaskCreationRequestInput';
import type { RecipePrepTaskUpdateRequestInput } from '../models/RecipePrepTaskUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipePrepTasksService {
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
}
