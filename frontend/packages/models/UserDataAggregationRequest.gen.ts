// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserDataAggregationRequest {
  reportID: string;
  userID: string;
}

export class UserDataAggregationRequest implements IUserDataAggregationRequest {
  reportID: string;
  userID: string;
  constructor(input: Partial<UserDataAggregationRequest> = {}) {
    this.reportID = input.reportID || '';
    this.userID = input.userID || '';
  }
}
