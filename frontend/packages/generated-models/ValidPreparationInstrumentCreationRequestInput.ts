// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidPreparationInstrumentCreationRequestInput {
  notes: string;
  validInstrumentID: string;
  validPreparationID: string;
}

export class ValidPreparationInstrumentCreationRequestInput implements IValidPreparationInstrumentCreationRequestInput {
  notes: string;
  validInstrumentID: string;
  validPreparationID: string;
  constructor(input: Partial<ValidPreparationInstrumentCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.validInstrumentID = input.validInstrumentID = '';
    this.validPreparationID = input.validPreparationID = '';
  }
}
