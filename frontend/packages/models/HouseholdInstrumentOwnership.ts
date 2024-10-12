// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidInstrument } from './ValidInstrument';


export interface IHouseholdInstrumentOwnership {
   archivedAt: string;
 belongsToHousehold: string;
 createdAt: string;
 id: string;
 instrument: ValidInstrument;
 lastUpdatedAt: string;
 notes: string;
 quantity: number;

}

export class HouseholdInstrumentOwnership implements IHouseholdInstrumentOwnership {
   archivedAt: string;
 belongsToHousehold: string;
 createdAt: string;
 id: string;
 instrument: ValidInstrument;
 lastUpdatedAt: string;
 notes: string;
 quantity: number;
constructor(input: Partial<HouseholdInstrumentOwnership> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.belongsToHousehold = input.belongsToHousehold || '';
 this.createdAt = input.createdAt || '';
 this.id = input.id || '';
 this.instrument = input.instrument || new ValidInstrument();
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.notes = input.notes || '';
 this.quantity = input.quantity || 0;
}
}