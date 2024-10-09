// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentCreationRequestInput {
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
}

export class ValidInstrumentCreationRequestInput implements IValidInstrumentCreationRequestInput {
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  constructor(input: Partial<ValidInstrumentCreationRequestInput> = {}) {
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
    this.description = input.description = '';
  }
}
