// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidPreparation } from './ValidPreparation';
 import { ValidVessel } from './ValidVessel';


export interface IValidPreparationVessel {
   instrument: ValidVessel;
 lastUpdatedAt?: string;
 notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;

}

export class ValidPreparationVessel implements IValidPreparationVessel {
   instrument: ValidVessel;
 lastUpdatedAt?: string;
 notes: string;
 preparation: ValidPreparation;
 archivedAt?: string;
 createdAt: string;
 id: string;
constructor(input: Partial<ValidPreparationVessel> = {}) {
	 this.instrument = input.instrument = new ValidVessel();
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.notes = input.notes = '';
 this.preparation = input.preparation = new ValidPreparation();
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
}
}