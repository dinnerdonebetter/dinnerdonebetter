// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentCreationRequestInput {
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
}

export class ValidInstrumentCreationRequestInput implements IValidInstrumentCreationRequestInput {
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  description: string;
  displayInSummaryLists: boolean;
  constructor(input: Partial<ValidInstrumentCreationRequestInput> = {}) {
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
    this.description = input.description = '';
    this.displayInSummaryLists = input.displayInSummaryLists = false;
  }
}
