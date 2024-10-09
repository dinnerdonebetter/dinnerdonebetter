// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInvitationUpdateRequestInput {
  token: string;
  note: string;
}

export class HouseholdInvitationUpdateRequestInput implements IHouseholdInvitationUpdateRequestInput {
  token: string;
  note: string;
  constructor(input: Partial<HouseholdInvitationUpdateRequestInput> = {}) {
    this.token = input.token = '';
    this.note = input.note = '';
  }
}
