// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdOwnershipTransferInput {
   reason: string;
 currentOwner: string;
 newOwner: string;

}

export class HouseholdOwnershipTransferInput implements IHouseholdOwnershipTransferInput {
   reason: string;
 currentOwner: string;
 newOwner: string;
constructor(input: Partial<HouseholdOwnershipTransferInput> = {}) {
	 this.reason = input.reason = '';
 this.currentOwner = input.currentOwner = '';
 this.newOwner = input.newOwner = '';
}
}