// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentUpdateRequestInput {
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
}

export class ValidInstrumentUpdateRequestInput implements IValidInstrumentUpdateRequestInput {
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  constructor(input: Partial<ValidInstrumentUpdateRequestInput> = {}) {
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
    this.usableForStorage = input.usableForStorage || false;
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.iconPath = input.iconPath || '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
  }
}
