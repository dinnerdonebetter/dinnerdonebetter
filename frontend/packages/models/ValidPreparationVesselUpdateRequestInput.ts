// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidPreparationVesselUpdateRequestInput {
  notes: string;
  validPreparationID: string;
  validVesselID: string;
}

export class ValidPreparationVesselUpdateRequestInput implements IValidPreparationVesselUpdateRequestInput {
  notes: string;
  validPreparationID: string;
  validVesselID: string;
  constructor(input: Partial<ValidPreparationVesselUpdateRequestInput> = {}) {
    this.notes = input.notes || '';
    this.validPreparationID = input.validPreparationID || '';
    this.validVesselID = input.validVesselID || '';
  }
}
