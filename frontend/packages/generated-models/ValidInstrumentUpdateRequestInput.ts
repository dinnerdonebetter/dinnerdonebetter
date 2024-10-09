// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentUpdateRequestInput {
  pluralName?: string;
  slug?: string;
  usableForStorage?: boolean;
  description?: string;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
}

export class ValidInstrumentUpdateRequestInput implements IValidInstrumentUpdateRequestInput {
  pluralName?: string;
  slug?: string;
  usableForStorage?: boolean;
  description?: string;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  constructor(input: Partial<ValidInstrumentUpdateRequestInput> = {}) {
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.usableForStorage = input.usableForStorage;
    this.description = input.description;
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.iconPath = input.iconPath;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.name = input.name;
  }
}
