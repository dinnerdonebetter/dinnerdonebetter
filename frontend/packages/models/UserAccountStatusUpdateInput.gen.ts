// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserAccountStatusUpdateInput {
  newStatus: string;
  reason: string;
  targetUserID: string;
}

export class UserAccountStatusUpdateInput implements IUserAccountStatusUpdateInput {
  newStatus: string;
  reason: string;
  targetUserID: string;
  constructor(input: Partial<UserAccountStatusUpdateInput> = {}) {
    this.newStatus = input.newStatus || '';
    this.reason = input.reason || '';
    this.targetUserID = input.targetUserID || '';
  }
}
