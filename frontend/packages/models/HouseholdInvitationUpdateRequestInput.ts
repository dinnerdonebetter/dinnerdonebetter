// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdInvitationUpdateRequestInput {
   note: string;
 token: string;

}

export class HouseholdInvitationUpdateRequestInput implements IHouseholdInvitationUpdateRequestInput {
   note: string;
 token: string;
constructor(input: Partial<HouseholdInvitationUpdateRequestInput> = {}) {
	 this.note = input.note || '';
 this.token = input.token || '';
}
}