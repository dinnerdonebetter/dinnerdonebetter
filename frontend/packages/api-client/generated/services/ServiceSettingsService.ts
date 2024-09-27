/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ServiceSetting } from '../models/ServiceSetting';
import type { ServiceSettingConfiguration } from '../models/ServiceSettingConfiguration';
import type { ServiceSettingConfigurationCreationRequestInput } from '../models/ServiceSettingConfigurationCreationRequestInput';
import type { ServiceSettingConfigurationUpdateRequestInput } from '../models/ServiceSettingConfigurationUpdateRequestInput';
import type { ServiceSettingCreationRequestInput } from '../models/ServiceSettingCreationRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ServiceSettingsService {
  /**
   * Operation for fetching ServiceSetting
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
  public static getServiceSettings(
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
      data?: Array<ServiceSetting>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings',
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
   * Operation for creating ServiceSetting
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createServiceSetting(requestBody: ServiceSettingCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ServiceSetting;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/settings',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ServiceSetting
   * @param serviceSettingId
   * @returns any
   * @throws ApiError
   */
  public static archiveServiceSetting(serviceSettingId: string): CancelablePromise<
    APIResponse & {
      data?: ServiceSetting;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/settings/{serviceSettingID}',
      path: {
        serviceSettingID: serviceSettingId,
      },
    });
  }
  /**
   * Operation for fetching ServiceSetting
   * @param serviceSettingId
   * @returns any
   * @throws ApiError
   */
  public static getServiceSetting(serviceSettingId: string): CancelablePromise<
    APIResponse & {
      data?: ServiceSetting;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings/{serviceSettingID}',
      path: {
        serviceSettingID: serviceSettingId,
      },
    });
  }
  /**
   * Operation for creating ServiceSettingConfiguration
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createServiceSettingConfiguration(
    requestBody: ServiceSettingConfigurationCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ServiceSettingConfiguration;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/settings/configurations',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ServiceSettingConfiguration
   * @param serviceSettingConfigurationId
   * @returns any
   * @throws ApiError
   */
  public static archiveServiceSettingConfiguration(serviceSettingConfigurationId: string): CancelablePromise<
    APIResponse & {
      data?: ServiceSettingConfiguration;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/settings/configurations/{serviceSettingConfigurationID}',
      path: {
        serviceSettingConfigurationID: serviceSettingConfigurationId,
      },
    });
  }
  /**
   * Operation for updating ServiceSettingConfiguration
   * @param serviceSettingConfigurationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateServiceSettingConfiguration(
    serviceSettingConfigurationId: string,
    requestBody: ServiceSettingConfigurationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ServiceSettingConfiguration;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/settings/configurations/{serviceSettingConfigurationID}',
      path: {
        serviceSettingConfigurationID: serviceSettingConfigurationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ServiceSettingConfiguration
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
  public static getServiceSettingConfigurationsForHousehold(
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
      data?: Array<ServiceSettingConfiguration>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings/configurations/household',
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
   * Operation for fetching ServiceSettingConfiguration
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
  public static getServiceSettingConfigurationsForUser(
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
      data?: Array<ServiceSettingConfiguration>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings/configurations/user',
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
   * Operation for fetching ServiceSettingConfiguration
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param serviceSettingConfigurationName
   * @returns any
   * @throws ApiError
   */
  public static getServiceSettingConfigurationByName(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    serviceSettingConfigurationName: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ServiceSettingConfiguration>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings/configurations/user/{serviceSettingConfigurationName}',
      path: {
        serviceSettingConfigurationName: serviceSettingConfigurationName,
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
   * Operation for fetching ServiceSetting
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
  public static searchForServiceSettings(
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
      data?: Array<ServiceSetting>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/settings/search',
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
