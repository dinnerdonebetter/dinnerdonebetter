/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DefaultService {
  /**
   * checks for service liveness
   * @throws ApiError
   */
  public static checkForLiveness(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/_meta_/live',
    });
  }
  /**
   * checks for service readiness
   * @throws ApiError
   */
  public static checkForReadiness(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/_meta_/ready',
    });
  }
}
