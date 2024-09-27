/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { UserNotification } from '../models/UserNotification';
import type { UserNotificationCreationRequestInput } from '../models/UserNotificationCreationRequestInput';
import type { UserNotificationUpdateRequestInput } from '../models/UserNotificationUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UserNotificationsService {
  /**
   * Operation for fetching UserNotification
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
  public static getUserNotifications(
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
      data?: Array<UserNotification>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/user_notifications',
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
   * Operation for creating UserNotification
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createUserNotification(requestBody: UserNotificationCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: UserNotification;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/user_notifications',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching UserNotification
   * @param userNotificationId
   * @returns any
   * @throws ApiError
   */
  public static getUserNotification(userNotificationId: string): CancelablePromise<
    APIResponse & {
      data?: UserNotification;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/user_notifications/{userNotificationID}',
      path: {
        userNotificationID: userNotificationId,
      },
    });
  }
  /**
   * Operation for updating UserNotification
   * @param userNotificationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateUserNotification(
    userNotificationId: string,
    requestBody: UserNotificationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: UserNotification;
    }
  > {
    return __request(OpenAPI, {
      method: 'PATCH',
      url: '/api/v1/user_notifications/{userNotificationID}',
      path: {
        userNotificationID: userNotificationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
