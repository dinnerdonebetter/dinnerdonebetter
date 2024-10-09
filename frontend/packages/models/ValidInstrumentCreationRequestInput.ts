// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentCreationRequestInput {
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
}

export class ValidInstrumentCreationRequestInput implements IValidInstrumentCreationRequestInput {
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  constructor(input: Partial<ValidInstrumentCreationRequestInput> = {}) {
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
    this.usableForStorage = input.usableForStorage || false;
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.iconPath = input.iconPath || '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.name = input.name || '';
  }
}
