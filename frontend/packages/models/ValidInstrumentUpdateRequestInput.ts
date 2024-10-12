// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidInstrumentUpdateRequestInput {
   description: string;
 displayInSummaryLists: boolean;
 iconPath: string;
 includeInGeneratedInstructions: boolean;
 name: string;
 pluralName: string;
 slug: string;
 usableForStorage: boolean;

}

export class ValidInstrumentUpdateRequestInput implements IValidInstrumentUpdateRequestInput {
   description: string;
 displayInSummaryLists: boolean;
 iconPath: string;
 includeInGeneratedInstructions: boolean;
 name: string;
 pluralName: string;
 slug: string;
 usableForStorage: boolean;
constructor(input: Partial<ValidInstrumentUpdateRequestInput> = {}) {
	 this.description = input.description || '';
 this.displayInSummaryLists = input.displayInSummaryLists || false;
 this.iconPath = input.iconPath || '';
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
 this.name = input.name || '';
 this.pluralName = input.pluralName || '';
 this.slug = input.slug || '';
 this.usableForStorage = input.usableForStorage || false;
}
}