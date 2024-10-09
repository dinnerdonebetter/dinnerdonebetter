// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentCreationRequestInput {
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
}

export class ValidInstrumentCreationRequestInput implements IValidInstrumentCreationRequestInput {
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  constructor(input: Partial<ValidInstrumentCreationRequestInput> = {}) {
    this.description = input.description = '';
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
  }
}
