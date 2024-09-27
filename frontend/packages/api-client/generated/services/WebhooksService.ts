/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { Webhook } from '../models/Webhook';
import type { WebhookCreationRequestInput } from '../models/WebhookCreationRequestInput';
import type { WebhookTriggerEvent } from '../models/WebhookTriggerEvent';
import type { WebhookTriggerEventCreationRequestInput } from '../models/WebhookTriggerEventCreationRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class WebhooksService {
  /**
   * Operation for fetching Webhook
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
  public static getWebhooks(
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
      data?: Array<Webhook>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/webhooks',
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
   * Operation for creating Webhook
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createWebhook(requestBody: WebhookCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: Webhook;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/webhooks',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving Webhook
   * @param webhookId
   * @returns any
   * @throws ApiError
   */
  public static archiveWebhook(webhookId: string): CancelablePromise<
    APIResponse & {
      data?: Webhook;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/webhooks/{webhookID}',
      path: {
        webhookID: webhookId,
      },
    });
  }
  /**
   * Operation for fetching Webhook
   * @param webhookId
   * @returns any
   * @throws ApiError
   */
  public static getWebhook(webhookId: string): CancelablePromise<
    APIResponse & {
      data?: Webhook;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/webhooks/{webhookID}',
      path: {
        webhookID: webhookId,
      },
    });
  }
  /**
   * Operation for creating WebhookTriggerEvent
   * @param webhookId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createWebhookTriggerEvent(
    webhookId: string,
    requestBody: WebhookTriggerEventCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: WebhookTriggerEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/webhooks/{webhookID}/trigger_events',
      path: {
        webhookID: webhookId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving WebhookTriggerEvent
   * @param webhookId
   * @param webhookTriggerEventId
   * @returns any
   * @throws ApiError
   */
  public static archiveWebhookTriggerEvent(
    webhookId: string,
    webhookTriggerEventId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: WebhookTriggerEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}',
      path: {
        webhookID: webhookId,
        webhookTriggerEventID: webhookTriggerEventId,
      },
    });
  }
}
