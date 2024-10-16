import { AccessToken } from 'simple-oauth2';
import { serialize } from 'cookie';

import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

export interface UserSessionDetails {
  token: AccessToken;
  userID: string;
  householdID: string;
}

export function cookieEncoderBuilder(cookieName: string, encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>) {
  return function (token: AccessToken, userID: string, householdID: string): string {
    const cookieValue = encryptorDecryptor.encrypt({
      token: token,
      userID: userID,
      householdID: householdID,
    });

    return serialize(cookieName, cookieValue, {
      path: '/',
      expires: new Date(JSON.parse(JSON.stringify(token))['expires_at']),
      httpOnly: true,
    });
  };
}
