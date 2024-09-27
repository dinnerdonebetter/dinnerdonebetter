/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidPreparation } from '../models/ValidPreparation';
import type { ValidPreparationCreationRequestInput } from '../models/ValidPreparationCreationRequestInput';
import type { ValidPreparationUpdateRequestInput } from '../models/ValidPreparationUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidPreparationsService {
  /**
   * Operation for fetching ValidPreparation
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
  public static getValidPreparations(
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
      data?: Array<ValidPreparation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparations',
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
   * Operation for creating ValidPreparation
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidPreparation(requestBody: ValidPreparationCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_preparations',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidPreparation
   * @param validPreparationId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidPreparation(validPreparationId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_preparations/{validPreparationID}',
      path: {
        validPreparationID: validPreparationId,
      },
    });
  }
  /**
   * Operation for fetching ValidPreparation
   * @param validPreparationId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparation(validPreparationId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparations/{validPreparationID}',
      path: {
        validPreparationID: validPreparationId,
      },
    });
  }
  /**
   * Operation for updating ValidPreparation
   * @param validPreparationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidPreparation(
    validPreparationId: string,
    requestBody: ValidPreparationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_preparations/{validPreparationID}',
      path: {
        validPreparationID: validPreparationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidPreparation
   * @returns any
   * @throws ApiError
   */
  public static getRandomValidPreparation(): CancelablePromise<
    APIResponse & {
      data?: ValidPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparations/random',
    });
  }
  /**
   * Operation for fetching ValidPreparation
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
  public static searchForValidPreparations(
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
      data?: Array<ValidPreparation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparations/search',
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
