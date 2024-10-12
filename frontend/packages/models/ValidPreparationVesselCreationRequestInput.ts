// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidPreparationVesselCreationRequestInput {
   notes: string;
 validPreparationID: string;
 validVesselID: string;

}

export class ValidPreparationVesselCreationRequestInput implements IValidPreparationVesselCreationRequestInput {
   notes: string;
 validPreparationID: string;
 validVesselID: string;
constructor(input: Partial<ValidPreparationVesselCreationRequestInput> = {}) {
	 this.notes = input.notes || '';
 this.validPreparationID = input.validPreparationID || '';
 this.validVesselID = input.validVesselID || '';
}
}