import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';
import { UserSessionDetails } from '@dinnerdonebetter/models';

export function cookieEncoderBuilder(encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>) {
  return function (jwt: string, userID: string, householdID: string): string {
    return encryptorDecryptor.encrypt({
      Token: jwt,
      UserID: userID,
      HouseholdID: householdID,
    });
  };
}
