// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdInstrumentOwnershipCreationRequestInput {
   notes: string;
 quantity: number;
 validInstrumentID: string;
 belongsToHousehold: string;

}

export class HouseholdInstrumentOwnershipCreationRequestInput implements IHouseholdInstrumentOwnershipCreationRequestInput {
   notes: string;
 quantity: number;
 validInstrumentID: string;
 belongsToHousehold: string;
constructor(input: Partial<HouseholdInstrumentOwnershipCreationRequestInput> = {}) {
	 this.notes = input.notes = '';
 this.quantity = input.quantity = 0;
 this.validInstrumentID = input.validInstrumentID = '';
 this.belongsToHousehold = input.belongsToHousehold = '';
}
}