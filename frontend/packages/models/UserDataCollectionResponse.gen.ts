// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserDataCollectionResponse {
  reportID: string;
}

export class UserDataCollectionResponse implements IUserDataCollectionResponse {
  reportID: string;
  constructor(input: Partial<UserDataCollectionResponse> = {}) {
    this.reportID = input.reportID || '';
  }
}
