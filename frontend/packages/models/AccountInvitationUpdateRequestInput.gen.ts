// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountInvitationUpdateRequestInput {
  note: string;
  token: string;
}

export class AccountInvitationUpdateRequestInput implements IAccountInvitationUpdateRequestInput {
  note: string;
  token: string;
  constructor(input: Partial<AccountInvitationUpdateRequestInput> = {}) {
    this.note = input.note || '';
    this.token = input.token || '';
  }
}
