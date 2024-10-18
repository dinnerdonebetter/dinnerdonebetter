// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInvitationCreationRequestInput {
  expiresAt: string;
  note: string;
  toEmail: string;
  toName: string;
}

export class HouseholdInvitationCreationRequestInput implements IHouseholdInvitationCreationRequestInput {
  expiresAt: string;
  note: string;
  toEmail: string;
  toName: string;
  constructor(input: Partial<HouseholdInvitationCreationRequestInput> = {}) {
    this.expiresAt = input.expiresAt || '';
    this.note = input.note || '';
    this.toEmail = input.toEmail || '';
    this.toName = input.toName || '';
  }
}
