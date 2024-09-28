/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { WebhookTriggerEvent } from '../models/WebhookTriggerEvent';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class WebhookTriggerEventsService {
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
