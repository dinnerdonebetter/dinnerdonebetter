// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidPreparationInstrumentUpdateRequestInput {
   notes: string;
 validInstrumentID: string;
 validPreparationID: string;

}

export class ValidPreparationInstrumentUpdateRequestInput implements IValidPreparationInstrumentUpdateRequestInput {
   notes: string;
 validInstrumentID: string;
 validPreparationID: string;
constructor(input: Partial<ValidPreparationInstrumentUpdateRequestInput> = {}) {
	 this.notes = input.notes || '';
 this.validInstrumentID = input.validInstrumentID || '';
 this.validPreparationID = input.validPreparationID || '';
}
}