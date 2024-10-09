// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdOwnershipTransferInput {
  currentOwner: string;
  newOwner: string;
  reason: string;
}

export class HouseholdOwnershipTransferInput implements IHouseholdOwnershipTransferInput {
  currentOwner: string;
  newOwner: string;
  reason: string;
  constructor(input: Partial<HouseholdOwnershipTransferInput> = {}) {
    this.currentOwner = input.currentOwner || '';
    this.newOwner = input.newOwner || '';
    this.reason = input.reason || '';
  }
}
