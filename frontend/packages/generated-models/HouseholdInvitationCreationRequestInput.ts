// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInvitationCreationRequestInput {
  toName: string;
  expiresAt?: string;
  note: string;
  toEmail: string;
}

export class HouseholdInvitationCreationRequestInput implements IHouseholdInvitationCreationRequestInput {
  toName: string;
  expiresAt?: string;
  note: string;
  toEmail: string;
  constructor(input: Partial<HouseholdInvitationCreationRequestInput> = {}) {
    this.toName = input.toName = '';
    this.expiresAt = input.expiresAt;
    this.note = input.note = '';
    this.toEmail = input.toEmail = '';
  }
}
