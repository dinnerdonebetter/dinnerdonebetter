/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { AvatarUpdateInput } from '../models/AvatarUpdateInput';
import type { ChangeActiveHouseholdInput } from '../models/ChangeActiveHouseholdInput';
import type { EmailAddressVerificationRequestInput } from '../models/EmailAddressVerificationRequestInput';
import type { Household } from '../models/Household';
import type { JWTResponse } from '../models/JWTResponse';
import type { PasswordResetToken } from '../models/PasswordResetToken';
import type { PasswordResetTokenCreationRequestInput } from '../models/PasswordResetTokenCreationRequestInput';
import type { PasswordResetTokenRedemptionRequestInput } from '../models/PasswordResetTokenRedemptionRequestInput';
import type { PasswordUpdateInput } from '../models/PasswordUpdateInput';
import type { TOTPSecretRefreshInput } from '../models/TOTPSecretRefreshInput';
import type { TOTPSecretRefreshResponse } from '../models/TOTPSecretRefreshResponse';
import type { TOTPSecretVerificationInput } from '../models/TOTPSecretVerificationInput';
import type { User } from '../models/User';
import type { UserAccountStatusUpdateInput } from '../models/UserAccountStatusUpdateInput';
import type { UserCreationResponse } from '../models/UserCreationResponse';
import type { UserDetailsUpdateRequestInput } from '../models/UserDetailsUpdateRequestInput';
import type { UserEmailAddressUpdateInput } from '../models/UserEmailAddressUpdateInput';
import type { UserLoginInput } from '../models/UserLoginInput';
import type { UsernameReminderRequestInput } from '../models/UsernameReminderRequestInput';
import type { UsernameUpdateInput } from '../models/UsernameUpdateInput';
import type { UserPermissionsRequestInput } from '../models/UserPermissionsRequestInput';
import type { UserPermissionsResponse } from '../models/UserPermissionsResponse';
import type { UserRegistrationInput } from '../models/UserRegistrationInput';
import type { UserStatusResponse } from '../models/UserStatusResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UsersService {
  /**
   * Operation for creating UserStatusResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static adminUpdateUserStatus(requestBody: UserAccountStatusUpdateInput): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/admin/users/status',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching User
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
  public static getUsers(
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
      data?: Array<User>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/users',
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
   * Operation for archiving User
   * @param userId
   * @returns any
   * @throws ApiError
   */
  public static archiveUser(userId: string): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/users/{userID}',
      path: {
        userID: userId,
      },
    });
  }
  /**
   * Operation for fetching User
   * @param userId
   * @returns any
   * @throws ApiError
   */
  public static getUser(userId: string): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/users/{userID}',
      path: {
        userID: userId,
      },
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static uploadUserAvatar(requestBody: AvatarUpdateInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/avatar/upload',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateUserDetails(requestBody: UserDetailsUpdateRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/users/details',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateUserEmailAddress(requestBody: UserEmailAddressUpdateInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/users/email_address',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static verifyUserEmailAddress(requestBody: EmailAddressVerificationRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/email_address_verification',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating Household
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static changeActiveHousehold(requestBody: ChangeActiveHouseholdInput): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/household/select',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating
   * @param requestBody
   * @throws ApiError
   */
  public static updatePassword(requestBody: PasswordUpdateInput): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/users/password/new',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating UserPermissionsResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static checkPermissions(requestBody: UserPermissionsRequestInput): CancelablePromise<
    APIResponse & {
      data?: UserPermissionsResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/permissions/check',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching User
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
  public static searchForUsers(
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
      data?: Array<User>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/users/search',
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
  /**
   * Operation for fetching User
   * @returns any
   * @throws ApiError
   */
  public static getSelf(): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/users/self',
    });
  }
  /**
   * Operation for creating TOTPSecretRefreshResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static refreshTotpSecret(requestBody: TOTPSecretRefreshInput): CancelablePromise<
    APIResponse & {
      data?: TOTPSecretRefreshResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/totp_secret/new',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateUserUsername(requestBody: UsernameUpdateInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/users/username',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Creates a new user
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createUser(requestBody: UserRegistrationInput): CancelablePromise<
    APIResponse & {
      data?: UserCreationResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static verifyEmailAddress(requestBody: EmailAddressVerificationRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/email_address/verify',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating UserStatusResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static login(requestBody: UserLoginInput): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/login',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating UserStatusResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static adminLogin(requestBody: UserLoginInput): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/login/admin',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating JWTResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static loginForJwt(requestBody: UserLoginInput): CancelablePromise<
    APIResponse & {
      data?: JWTResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/login/jwt',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating JWTResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static adminLoginForJwt(requestBody: UserLoginInput): CancelablePromise<
    APIResponse & {
      data?: JWTResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/login/jwt/admin',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating UserStatusResponse
   * @returns any
   * @throws ApiError
   */
  public static logout(): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/logout',
    });
  }
  /**
   * Operation for creating PasswordResetToken
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static requestPasswordResetToken(requestBody: PasswordResetTokenCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: PasswordResetToken;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/password/reset',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for redeeming a password reset token
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static redeemPasswordResetToken(requestBody: PasswordResetTokenRedemptionRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/password/reset/redeem',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static verifyTotpSecret(requestBody: TOTPSecretVerificationInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/totp_secret/verify',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static requestUsernameReminder(requestBody: UsernameReminderRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/username/reminder',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
