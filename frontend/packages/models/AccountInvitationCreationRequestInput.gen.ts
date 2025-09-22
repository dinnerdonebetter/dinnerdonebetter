// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountInvitationCreationRequestInput {
  expiresAt: string;
  note: string;
  toEmail: string;
  toName: string;
}

export class AccountInvitationCreationRequestInput implements IAccountInvitationCreationRequestInput {
  expiresAt: string;
  note: string;
  toEmail: string;
  toName: string;
  constructor(input: Partial<AccountInvitationCreationRequestInput> = {}) {
    this.expiresAt = input.expiresAt || '';
    this.note = input.note || '';
    this.toEmail = input.toEmail || '';
    this.toName = input.toName || '';
  }
}
