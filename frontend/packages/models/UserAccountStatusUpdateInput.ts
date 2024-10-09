// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserAccountStatusUpdateInput {
  targetUserID: string;
  newStatus: string;
  reason: string;
}

export class UserAccountStatusUpdateInput implements IUserAccountStatusUpdateInput {
  targetUserID: string;
  newStatus: string;
  reason: string;
  constructor(input: Partial<UserAccountStatusUpdateInput> = {}) {
    this.targetUserID = input.targetUserID || '';
    this.newStatus = input.newStatus || '';
    this.reason = input.reason || '';
  }
}
