/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { UserIngredientPreference } from '../models/UserIngredientPreference';
import type { UserIngredientPreferenceCreationRequestInput } from '../models/UserIngredientPreferenceCreationRequestInput';
import type { UserIngredientPreferenceUpdateRequestInput } from '../models/UserIngredientPreferenceUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UserIngredientPreferencesService {
  /**
   * Operation for fetching UserIngredientPreference
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
  public static getUserIngredientPreferences(
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
      data?: Array<UserIngredientPreference>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/user_ingredient_preferences',
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
   * Operation for creating UserIngredientPreference
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createUserIngredientPreference(
    requestBody: UserIngredientPreferenceCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: UserIngredientPreference;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/user_ingredient_preferences',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving UserIngredientPreference
   * @param userIngredientPreferenceId
   * @returns any
   * @throws ApiError
   */
  public static archiveUserIngredientPreference(userIngredientPreferenceId: string): CancelablePromise<
    APIResponse & {
      data?: UserIngredientPreference;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/user_ingredient_preferences/{userIngredientPreferenceID}',
      path: {
        userIngredientPreferenceID: userIngredientPreferenceId,
      },
    });
  }
  /**
   * Operation for updating UserIngredientPreference
   * @param userIngredientPreferenceId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateUserIngredientPreference(
    userIngredientPreferenceId: string,
    requestBody: UserIngredientPreferenceUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: UserIngredientPreference;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/user_ingredient_preferences/{userIngredientPreferenceID}',
      path: {
        userIngredientPreferenceID: userIngredientPreferenceId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
