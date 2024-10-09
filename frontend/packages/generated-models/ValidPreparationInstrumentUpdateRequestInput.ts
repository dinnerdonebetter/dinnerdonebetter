// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidPreparationInstrumentUpdateRequestInput {
   validInstrumentID?: string;
 validPreparationID?: string;
 notes?: string;

}

export class ValidPreparationInstrumentUpdateRequestInput implements IValidPreparationInstrumentUpdateRequestInput {
   validInstrumentID?: string;
 validPreparationID?: string;
 notes?: string;
constructor(input: Partial<ValidPreparationInstrumentUpdateRequestInput> = {}) {
	 this.validInstrumentID = input.validInstrumentID;
 this.validPreparationID = input.validPreparationID;
 this.notes = input.notes;
}
}