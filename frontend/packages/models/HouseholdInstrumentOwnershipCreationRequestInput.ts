// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdInstrumentOwnershipCreationRequestInput {
   belongsToHousehold: string;
 notes: string;
 quantity: number;
 validInstrumentID: string;

}

export class HouseholdInstrumentOwnershipCreationRequestInput implements IHouseholdInstrumentOwnershipCreationRequestInput {
   belongsToHousehold: string;
 notes: string;
 quantity: number;
 validInstrumentID: string;
constructor(input: Partial<HouseholdInstrumentOwnershipCreationRequestInput> = {}) {
	 this.belongsToHousehold = input.belongsToHousehold || '';
 this.notes = input.notes || '';
 this.quantity = input.quantity || 0;
 this.validInstrumentID = input.validInstrumentID || '';
}
}