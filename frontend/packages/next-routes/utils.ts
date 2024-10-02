import { AccessToken } from 'simple-oauth2';
import { serialize } from 'cookie';

import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

export interface UserSessionDetails {
  Token: AccessToken;
  UserID: string;
  HouseholdID: string;
}

export function cookieEncoderBuilder(cookieName: string, encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>) {
  return function (token: AccessToken, userID: string, householdID: string): string {
    const cookieValue = encryptorDecryptor.encrypt({
      Token: token,
      UserID: userID,
      HouseholdID: householdID,
    });

    return serialize(cookieName, cookieValue, {
      path: '/',
      expires: new Date(new Date().getTime() + 24 * 60 * 60 * 1000),
      httpOnly: true,
    });
  };
}
