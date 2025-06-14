// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountOwnershipTransferInput {
  currentOwner: string;
  newOwner: string;
  reason: string;
}

export class AccountOwnershipTransferInput implements IAccountOwnershipTransferInput {
  currentOwner: string;
  newOwner: string;
  reason: string;
  constructor(input: Partial<AccountOwnershipTransferInput> = {}) {
    this.currentOwner = input.currentOwner || '';
    this.newOwner = input.newOwner || '';
    this.reason = input.reason || '';
  }
}
