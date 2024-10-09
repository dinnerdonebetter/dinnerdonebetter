// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserAccountStatusUpdateInput {
  reason: string;
  targetUserID: string;
  newStatus: string;
}

export class UserAccountStatusUpdateInput implements IUserAccountStatusUpdateInput {
  reason: string;
  targetUserID: string;
  newStatus: string;
  constructor(input: Partial<UserAccountStatusUpdateInput> = {}) {
    this.reason = input.reason = '';
    this.targetUserID = input.targetUserID = '';
    this.newStatus = input.newStatus = '';
  }
}
