// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidInstrument } from './ValidInstrument';
 import { ValidPreparation } from './ValidPreparation';


export interface IValidPreparationInstrument {
   lastUpdatedAt?: string;
 notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;
 instrument: ValidInstrument;

}

export class ValidPreparationInstrument implements IValidPreparationInstrument {
   lastUpdatedAt?: string;
 notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;
 instrument: ValidInstrument;
constructor(input: Partial<ValidPreparationInstrument> = {}) {
	 this.lastUpdatedAt = input.lastUpdatedAt;
 this.notes = input.notes = '';
 this.preparation = input.preparation = new ValidPreparation();
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.instrument = input.instrument = new ValidInstrument();
}
}