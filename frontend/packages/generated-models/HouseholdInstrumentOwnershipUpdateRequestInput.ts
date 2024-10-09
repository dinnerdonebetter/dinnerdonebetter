// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInstrumentOwnershipUpdateRequestInput {
  notes?: string;
  quantity?: number;
  validInstrumentID?: string;
}

export class HouseholdInstrumentOwnershipUpdateRequestInput implements IHouseholdInstrumentOwnershipUpdateRequestInput {
  notes?: string;
  quantity?: number;
  validInstrumentID?: string;
  constructor(input: Partial<HouseholdInstrumentOwnershipUpdateRequestInput> = {}) {
    this.notes = input.notes;
    this.quantity = input.quantity;
    this.validInstrumentID = input.validInstrumentID;
  }
}
