/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { MealPlanEvent } from '../models/MealPlanEvent';
import type { MealPlanEventCreationRequestInput } from '../models/MealPlanEventCreationRequestInput';
import type { MealPlanEventUpdateRequestInput } from '../models/MealPlanEventUpdateRequestInput';
import type { MealPlanOption } from '../models/MealPlanOption';
import type { MealPlanOptionCreationRequestInput } from '../models/MealPlanOptionCreationRequestInput';
import type { MealPlanOptionUpdateRequestInput } from '../models/MealPlanOptionUpdateRequestInput';
import type { MealPlanOptionVote } from '../models/MealPlanOptionVote';
import type { MealPlanOptionVoteCreationRequestInput } from '../models/MealPlanOptionVoteCreationRequestInput';
import type { MealPlanOptionVoteUpdateRequestInput } from '../models/MealPlanOptionVoteUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class MealPlanEventsService {
  /**
   * Operation for fetching MealPlanEvent
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param mealPlanId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanEvents(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    mealPlanId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<MealPlanEvent>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events',
      path: {
        mealPlanID: mealPlanId,
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
   * Operation for creating MealPlanEvent
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanEvent(
    mealPlanId: string,
    requestBody: MealPlanEventCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/events',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving MealPlanEvent
   * @param mealPlanId
   * @param mealPlanEventId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlanEvent(
    mealPlanId: string,
    mealPlanEventId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanEvent
   * @param mealPlanId
   * @param mealPlanEventId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanEvent(
    mealPlanId: string,
    mealPlanEventId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
      },
    });
  }
  /**
   * Operation for updating MealPlanEvent
   * @param mealPlanId
   * @param mealPlanEventId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanEvent(
    mealPlanId: string,
    mealPlanEventId: string,
    requestBody: MealPlanEventUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanEvent;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching MealPlanOption
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param mealPlanId
   * @param mealPlanEventId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanOptions(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    mealPlanId: string,
    mealPlanEventId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<MealPlanOption>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
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
   * Operation for creating MealPlanOption
   * @param mealPlanId
   * @param mealPlanEventId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanOption(
    mealPlanId: string,
    mealPlanEventId: string,
    requestBody: MealPlanOptionCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOption;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving MealPlanOption
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlanOption(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOption;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanOption
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanOption(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOption;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
      },
    });
  }
  /**
   * Operation for updating MealPlanOption
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanOption(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
    requestBody: MealPlanOptionUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOption;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching MealPlanOptionVote
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanOptionVotes(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<MealPlanOptionVote>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
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
   * Operation for archiving MealPlanOptionVote
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @param mealPlanOptionVoteId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlanOptionVote(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
    mealPlanOptionVoteId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOptionVote;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
        mealPlanOptionVoteID: mealPlanOptionVoteId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanOptionVote
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @param mealPlanOptionVoteId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanOptionVote(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
    mealPlanOptionVoteId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOptionVote;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
        mealPlanOptionVoteID: mealPlanOptionVoteId,
      },
    });
  }
  /**
   * Operation for updating MealPlanOptionVote
   * @param mealPlanId
   * @param mealPlanEventId
   * @param mealPlanOptionId
   * @param mealPlanOptionVoteId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanOptionVote(
    mealPlanId: string,
    mealPlanEventId: string,
    mealPlanOptionId: string,
    mealPlanOptionVoteId: string,
    requestBody: MealPlanOptionVoteUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOptionVote;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
        mealPlanOptionID: mealPlanOptionId,
        mealPlanOptionVoteID: mealPlanOptionVoteId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating MealPlanOptionVote
   * @param mealPlanId
   * @param mealPlanEventId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanOptionVote(
    mealPlanId: string,
    mealPlanEventId: string,
    requestBody: MealPlanOptionVoteCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanOptionVote;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote',
      path: {
        mealPlanID: mealPlanId,
        mealPlanEventID: mealPlanEventId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
