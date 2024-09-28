/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidMeasurementUnitConversion } from '../models/ValidMeasurementUnitConversion';
import type { ValidMeasurementUnitConversionCreationRequestInput } from '../models/ValidMeasurementUnitConversionCreationRequestInput';
import type { ValidMeasurementUnitConversionUpdateRequestInput } from '../models/ValidMeasurementUnitConversionUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidMeasurementConversionsService {
  /**
   * Operation for creating ValidMeasurementUnitConversion
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidMeasurementUnitConversion(
    requestBody: ValidMeasurementUnitConversionCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_measurement_conversions',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidMeasurementUnitConversion
   * @param validMeasurementUnitConversionId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidMeasurementUnitConversion(validMeasurementUnitConversionId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}',
      path: {
        validMeasurementUnitConversionID: validMeasurementUnitConversionId,
      },
    });
  }
  /**
   * Operation for fetching ValidMeasurementUnitConversion
   * @param validMeasurementUnitConversionId
   * @returns any
   * @throws ApiError
   */
  public static getValidMeasurementUnitConversion(validMeasurementUnitConversionId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}',
      path: {
        validMeasurementUnitConversionID: validMeasurementUnitConversionId,
      },
    });
  }
  /**
   * Operation for updating ValidMeasurementUnitConversion
   * @param validMeasurementUnitConversionId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidMeasurementUnitConversion(
    validMeasurementUnitConversionId: string,
    requestBody: ValidMeasurementUnitConversionUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}',
      path: {
        validMeasurementUnitConversionID: validMeasurementUnitConversionId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidMeasurementUnitConversion
   * @param validMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
      },
    });
  }
  /**
   * Operation for fetching ValidMeasurementUnitConversion
   * @param validMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static validMeasurementUnitConversionsToUnit(validMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnitConversion;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
      },
    });
  }
}
