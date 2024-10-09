// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInvitationCreationRequestInput {
  toEmail: string;
  toName: string;
  expiresAt: string;
  note: string;
}

export class HouseholdInvitationCreationRequestInput implements IHouseholdInvitationCreationRequestInput {
  toEmail: string;
  toName: string;
  expiresAt: string;
  note: string;
  constructor(input: Partial<HouseholdInvitationCreationRequestInput> = {}) {
    this.toEmail = input.toEmail || '';
    this.toName = input.toName || '';
    this.expiresAt = input.expiresAt || '';
    this.note = input.note || '';
  }
}
