/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidInstrument } from '../models/ValidInstrument';
import type { ValidInstrumentCreationRequestInput } from '../models/ValidInstrumentCreationRequestInput';
import type { ValidInstrumentUpdateRequestInput } from '../models/ValidInstrumentUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidInstrumentsService {
  /**
   * Operation for fetching ValidInstrument
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
  public static getValidInstruments(
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
      data?: Array<ValidInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_instruments',
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
   * Operation for creating ValidInstrument
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidInstrument(requestBody: ValidInstrumentCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_instruments',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidInstrument
   * @param validInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidInstrument(validInstrumentId: string): CancelablePromise<
    APIResponse & {
      data?: ValidInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_instruments/{validInstrumentID}',
      path: {
        validInstrumentID: validInstrumentId,
      },
    });
  }
  /**
   * Operation for fetching ValidInstrument
   * @param validInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static getValidInstrument(validInstrumentId: string): CancelablePromise<
    APIResponse & {
      data?: ValidInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_instruments/{validInstrumentID}',
      path: {
        validInstrumentID: validInstrumentId,
      },
    });
  }
  /**
   * Operation for updating ValidInstrument
   * @param validInstrumentId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidInstrument(
    validInstrumentId: string,
    requestBody: ValidInstrumentUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_instruments/{validInstrumentID}',
      path: {
        validInstrumentID: validInstrumentId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidInstrument
   * @returns any
   * @throws ApiError
   */
  public static getRandomValidInstrument(): CancelablePromise<
    APIResponse & {
      data?: ValidInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_instruments/random',
    });
  }
  /**
   * Operation for fetching ValidInstrument
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
  public static searchForValidInstruments(
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
      data?: Array<ValidInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_instruments/search',
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
