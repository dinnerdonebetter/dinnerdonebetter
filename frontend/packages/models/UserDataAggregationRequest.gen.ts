// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserDataAggregationRequest {
  id: string;
  reportID: string;
  userID: string;
}

export class UserDataAggregationRequest implements IUserDataAggregationRequest {
  id: string;
  reportID: string;
  userID: string;
  constructor(input: Partial<UserDataAggregationRequest> = {}) {
    this.id = input.id || '';
    this.reportID = input.reportID || '';
    this.userID = input.userID || '';
  }
}
