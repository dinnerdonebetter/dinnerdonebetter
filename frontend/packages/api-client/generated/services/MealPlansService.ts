/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { FinalizeMealPlansResponse } from '../models/FinalizeMealPlansResponse';
import type { MealPlan } from '../models/MealPlan';
import type { MealPlanCreationRequestInput } from '../models/MealPlanCreationRequestInput';
import type { MealPlanEvent } from '../models/MealPlanEvent';
import type { MealPlanEventCreationRequestInput } from '../models/MealPlanEventCreationRequestInput';
import type { MealPlanEventUpdateRequestInput } from '../models/MealPlanEventUpdateRequestInput';
import type { MealPlanGroceryListItem } from '../models/MealPlanGroceryListItem';
import type { MealPlanGroceryListItemCreationRequestInput } from '../models/MealPlanGroceryListItemCreationRequestInput';
import type { MealPlanGroceryListItemUpdateRequestInput } from '../models/MealPlanGroceryListItemUpdateRequestInput';
import type { MealPlanOption } from '../models/MealPlanOption';
import type { MealPlanOptionCreationRequestInput } from '../models/MealPlanOptionCreationRequestInput';
import type { MealPlanOptionUpdateRequestInput } from '../models/MealPlanOptionUpdateRequestInput';
import type { MealPlanOptionVote } from '../models/MealPlanOptionVote';
import type { MealPlanOptionVoteCreationRequestInput } from '../models/MealPlanOptionVoteCreationRequestInput';
import type { MealPlanOptionVoteUpdateRequestInput } from '../models/MealPlanOptionVoteUpdateRequestInput';
import type { MealPlanTask } from '../models/MealPlanTask';
import type { MealPlanTaskCreationRequestInput } from '../models/MealPlanTaskCreationRequestInput';
import type { MealPlanTaskStatusChangeRequestInput } from '../models/MealPlanTaskStatusChangeRequestInput';
import type { MealPlanUpdateRequestInput } from '../models/MealPlanUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class MealPlansService {
  /**
   * Operation for fetching MealPlan
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
  public static getMealPlans(
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
      data?: Array<MealPlan>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans',
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
   * Operation for creating MealPlan
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlan(requestBody: MealPlanCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: MealPlan;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving MealPlan
   * @param mealPlanId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlan(mealPlanId: string): CancelablePromise<
    APIResponse & {
      data?: MealPlan;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}',
      path: {
        mealPlanID: mealPlanId,
      },
    });
  }
  /**
   * Operation for fetching MealPlan
   * @param mealPlanId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlan(mealPlanId: string): CancelablePromise<
    APIResponse & {
      data?: MealPlan;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}',
      path: {
        mealPlanID: mealPlanId,
      },
    });
  }
  /**
   * Operation for updating MealPlan
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlan(
    mealPlanId: string,
    requestBody: MealPlanUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlan;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
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
  /**
   * Operation for creating FinalizeMealPlansResponse
   * @param mealPlanId
   * @returns any
   * @throws ApiError
   */
  public static finalizeMealPlan(mealPlanId: string): CancelablePromise<
    APIResponse & {
      data?: FinalizeMealPlansResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/finalize',
      path: {
        mealPlanID: mealPlanId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanGroceryListItem
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
  public static getMealPlanGroceryListItemsForMealPlan(
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
      data?: Array<MealPlanGroceryListItem>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items',
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
   * Operation for creating MealPlanGroceryListItem
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanGroceryListItem(
    mealPlanId: string,
    requestBody: MealPlanGroceryListItemCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
    });
  }
  /**
   * Operation for updating MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
    requestBody: MealPlanGroceryListItemUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching MealPlanTask
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
  public static getMealPlanTasks(
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
      data?: Array<MealPlanTask>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks',
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
   * Operation for creating MealPlanTask
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanTask(
    mealPlanId: string,
    requestBody: MealPlanTaskCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching MealPlanTask
   * @param mealPlanId
   * @param mealPlanTaskId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanTask(
    mealPlanId: string,
    mealPlanTaskId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanTaskID: mealPlanTaskId,
      },
    });
  }
  /**
   * Operation for updating MealPlanTask
   * @param mealPlanId
   * @param mealPlanTaskId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanTaskStatus(
    mealPlanId: string,
    mealPlanTaskId: string,
    requestBody: MealPlanTaskStatusChangeRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'PATCH',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanTaskID: mealPlanTaskId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
