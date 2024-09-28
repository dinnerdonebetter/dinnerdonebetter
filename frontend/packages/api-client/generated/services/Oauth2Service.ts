/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { OAuth2Client } from '../models/OAuth2Client';
import type { OAuth2ClientCreationRequestInput } from '../models/OAuth2ClientCreationRequestInput';
import type { OAuth2ClientCreationResponse } from '../models/OAuth2ClientCreationResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class Oauth2Service {
  /**
   * Operation for fetching OAuth2Client
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
  public static getOAuth2Clients(
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
      data?: Array<OAuth2Client>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/oauth2_clients',
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
   * Operation for creating OAuth2ClientCreationResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createOAuth2Client(requestBody: OAuth2ClientCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: OAuth2ClientCreationResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/oauth2_clients',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving OAuth2Client
   * @param oauth2ClientId
   * @returns any
   * @throws ApiError
   */
  public static archiveOAuth2Client(oauth2ClientId: string): CancelablePromise<
    APIResponse & {
      data?: OAuth2Client;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/oauth2_clients/{oauth2ClientID}',
      path: {
        oauth2ClientID: oauth2ClientId,
      },
    });
  }
  /**
   * Operation for fetching OAuth2Client
   * @param oauth2ClientId
   * @returns any
   * @throws ApiError
   */
  public static getOAuth2Client(oauth2ClientId: string): CancelablePromise<
    APIResponse & {
      data?: OAuth2Client;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/oauth2_clients/{oauth2ClientID}',
      path: {
        oauth2ClientID: oauth2ClientId,
      },
    });
  }
  /**
   * Operation for fetching
   * @throws ApiError
   */
  public static getOauth2Authorize(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/oauth2/authorize',
    });
  }
  /**
   * Operation for creating
   * @throws ApiError
   */
  public static postOauth2Token(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/oauth2/token',
    });
  }
}
