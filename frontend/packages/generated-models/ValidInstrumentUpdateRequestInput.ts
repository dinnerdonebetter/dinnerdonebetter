// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentUpdateRequestInput {
  usableForStorage?: boolean;
  description?: string;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
}

export class ValidInstrumentUpdateRequestInput implements IValidInstrumentUpdateRequestInput {
  usableForStorage?: boolean;
  description?: string;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  constructor(input: Partial<ValidInstrumentUpdateRequestInput> = {}) {
    this.usableForStorage = input.usableForStorage;
    this.description = input.description;
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.iconPath = input.iconPath;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
  }
}
